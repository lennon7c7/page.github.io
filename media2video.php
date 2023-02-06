#!/usr/bin/env php
<?php
/**
 * media to video
 */
set_time_limit(0);

include './include/common.php';

$input_file = getInputImg();
foreach ($input_file as $dir => $value) {
    foreach ($value as $dir2 => $files) {
        $video_output_path = "./video/" . basename($dir);
        if (!is_dir($video_output_path)) {
            $res = mkdir($video_output_path, 0777, true);
            if (!$res) {
                echo "mkdir $video_output_path error" . PHP_EOL;
                return;
            }
        }

        $video_output_path .= '/' . basename($dir2);

        list($max_width, $max_height) = getPCRectangle($files);
        serializeJpgFilename($dir2);
        $input_img_template = $dir2 . '/%1d.jpg';
        $video_output = "$video_output_path.mp4";
        if (!file_exists($video_output)) {
            // -framerate 1/2 每张图显示2s
            // -r 30 30帧/秒
            // scale 把原图修改下分辨率，缺少的地方不剪切不拉伸而是加黑边，再把所有处理后的图片二次处理成视频
            $shell = "ffmpeg -framerate 1/2 -start_number 1 -i \"$input_img_template\" -c:v libx264 -r 30 -vf \"scale=" . $max_width . ':' . $max_height . ':force_original_aspect_ratio=decrease,pad=' . $max_width . ':' . $max_height . ':(ow-iw)/2:(oh-ih)/2" -qscale 1 "' . $video_output . '" -y';
            $out = [];
            exec($shell, $out);
        }


        $music_output = "$video_output_path-music.mp4";
        if (!file_exists($music_output)) {
            $audio_input_file = getRandomAudioFile();
            $shell = "ffmpeg -i \"$video_output\" -stream_loop -1 -i \"$audio_input_file\" -shortest -map 0:v:0 -map 1:a:0 -c:v copy \"$music_output\" -y";
            $out = [];
            exec($shell, $out);
        }
    }
}
