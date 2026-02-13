package core

import (
	"context"
	"errors"
	"io"
	"log"
	"os"
	"time"

	"github.com/weihaoli/szem/pkg/repo"
	"github.com/weihaoli/szem/pkg/storage"
)

type Worker struct {
	Jobs    repo.JobRepo
	Objects repo.ObjectRepo
	Storage storage.Storage
}

func (w *Worker) Start() {
	go w.loop()
}

func (w *Worker) loop() {
	for {
		time.Sleep(2 * time.Second)

		ctx := context.Background()

		job, err := w.Jobs.FetchOnePending(ctx)
		if err != nil {
			// 没有 pending job 或 DB 暂时错误：都继续下一轮
			continue
		}

		_ = w.Jobs.MarkRunning(ctx, job.ID)

		obj, err := w.Objects.GetByID(ctx, job.ObjectID)
		if err != nil {
			_ = w.Jobs.MarkFailed(ctx, job.ID, "object not found for job")
			continue
		}

		if obj.StagingPath == "" {
			_ = w.Jobs.MarkFailed(ctx, job.ID, "empty staging path")
			_ = w.Objects.MarkFailed(ctx, obj.ObjectID, "empty staging path")
			continue
		}

		f, err := os.Open(obj.StagingPath)
		if err != nil {
			_ = w.Jobs.MarkFailed(ctx, job.ID, "open staging failed: "+err.Error())
			_ = w.Objects.MarkFailed(ctx, obj.ObjectID, "open staging failed: "+err.Error())
			continue
		}

		err = w.Storage.Store(ctx, obj.Bucket, obj.Key, f, obj.Size)
		_ = f.Close()

		if err != nil {
			_ = w.Jobs.MarkFailed(ctx, job.ID, "storage store failed: "+err.Error())
			_ = w.Objects.MarkFailed(ctx, obj.ObjectID, "storage store failed: "+err.Error())
			continue
		}

		// 标记可用
		if err := w.Objects.MarkAvailable(ctx, obj.ObjectID); err != nil {
			_ = w.Jobs.MarkFailed(ctx, job.ID, "mark available failed: "+err.Error())
			continue
		}

		_ = w.Jobs.MarkDone(ctx, job.ID)

		// 成功后删除 staging
		if err := os.Remove(obj.StagingPath); err != nil && !errors.Is(err, os.ErrNotExist) {
			log.Println("warn: remove staging failed:", err)
		}
	}
}

// 小工具：把 reader 写到文件（PUT 阶段用）
func WriteToFile(path string, r io.Reader) error {
	if err := os.MkdirAll(dirOf(path), 0755); err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, r)
	return err
}

func dirOf(p string) string {
	// 不引入 filepath，避免你误以为必须；当然也可以用 filepath.Dir
	for i := len(p) - 1; i >= 0; i-- {
		if p[i] == '/' {
			return p[:i]
		}
	}
	return "."
}

