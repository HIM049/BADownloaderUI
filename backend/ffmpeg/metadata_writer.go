package ffmpeg

import (
	"bili-audio-downloader/backend/constants"
	"bili-audio-downloader/backend/utils"
	"fmt"
)

func WriteMetadata(input, output, coverPath, songName, songAuthor string, format string) error {
	err := initFfmpeg(&input, &output)
	if err != nil {
		return err
	}

	// 基础 ffmpeg 参数
	args := []string{
		"-y",
		"-i", input,
	}

	// 加入封面图（如果提供）
	if coverPath != "" {
		args = append(args, "-i", coverPath)

	}

	// 添加元数据字段
	args = append(args,
		"-map", "0", // 主音频流
	)

	if coverPath != "" {
		args = append(args,
			"-map", "1", // 图片作为第二个输入
			"-disposition:v", "attached_pic", // 标记为封面图
		)
	}

	// 元数据字段设置
	metadata := []string{
		"-metadata", "title=" + songName, // 标题
		"-metadata", "artist=" + songAuthor, // 艺术家
		//"-metadata", "album=Album Name", // 专辑名称
		//"-metadata", "album_artist=Album Artist", // 专辑艺术家
		//"-metadata", "composer=Composer Name", // 作曲者
		//"-metadata", "genre=Pop", // 音乐风格
		//"-metadata", "track=2/10", // 曲目号/总数
		//"-metadata", "comment=Demo release", // 注释
		//"-metadata", "lyrics=These are the lyrics...", // 歌词
		//"-metadata", "encoded_by=ffmpeg 6.0", // 编码者
		//"-metadata", "language=eng", // 语言
		//"-metadata", "grouping=Electronic", // 分组（可用于分类）
	}

	args = append(args, metadata...)

	// 编码设置：不转码音频，只复制流
	args = append(args, "-c", "copy")

	// 如果是 mp3，推荐设置 id3v2 版本
	if format == constants.AudioType.Mp3 {
		args = append(args, "-id3v2_version", "3")
	}

	args = append(args, "-f", format[1:], output)

	log, err := utils.RunCommand("ffmpeg", args...)
	if err != nil {
		return fmt.Errorf("ffmpeg error: %v\nOutput: %s", err, log)
	}

	return nil
}
