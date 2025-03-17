package server

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"gocaptcha/pkg/log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Run 运行 http 服务器
func Run(handler http.Handler, addr string) {
	srv := &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadHeaderTimeout: 2 * time.Second,
	}

	// 启动服务器协程
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.L().Fatal("服务遇到了异常", zap.Error(err))
		}
	}()

	// 阻塞并监听结束信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.L().Info("正在关闭服务...")

	// 关闭服务器（5秒超时时间）
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.L().Error("服务关闭失败", zap.Error(err))
	}

	log.L().Info("服务已关闭")
}
