#!/usr/bin/env php
<?php
/**
 * media to video
 */
set_time_limit(0);

include './include/common.php';

$input_file = getInputImg();
foreach ($input_file as $key => $value1) {
    foreach ($value1 as $key2 => $files) {
        $video_output_path = "./video/$key/";
        if (!is_dir($video_output_path)) {
            $res = mkdir($video_output_path, 0777, true);
            if (!$res) {
                echo "mkdir $video_output_path error" . PHP_EOL;
                return;
            }
        }


        $temp_file1 = [];
        foreach ($files as $filename) {
            if (!isImgFileExt($filename)) {
                continue;
            }

            $new_filename = covertImage($filename);
            if (empty($new_filename)) {
                continue;
            }

            $temp_file1[] = $new_filename;
        }

        if (empty($temp_file1)) {
            return;
        }

        $input_img_template = dirname($temp_file1[0]) . '/temp-' . date('Ymd') . '-%1d.jpg';
        $file_list = [];
        foreach ($temp_file1 as $key => $filename) {
            $temp_filename = str_replace('%1d', $key + 1, $input_img_template);;
            list($width, $height, $type, $attr) = getimagesize($filename);

            $width_not_divisible_by_2 = false;
            if ($width % 2 != 0) {
                $width_not_divisible_by_2 = true;
            }

            $height_not_divisible_by_2 = false;
            if ($height % 2 != 0) {
                $height_not_divisible_by_2 = true;
            }

            if ($width_not_divisible_by_2 || $height_not_divisible_by_2) {
                if ($width_not_divisible_by_2) {
                    $width--;
                }

                if ($height_not_divisible_by_2) {
                    $height--;
                }

                resizeImage($filename, $width, $height);
            }

            $file_list[] = [
                'old' => $filename,
                'new' => $temp_filename,
                'width' => $width,
                'height' => $height,
            ];
        }


        foreach ($file_list as $filename) {
            rename($filename['old'], $filename['new']);
        }


        list($max_width, $max_height) = getPCRectangle($file_list);
        $video_output = "$video_output_path$key2.mp4";
        if (!file_exists($video_output)) {
            // -framerate 1/2 每张图显示2s
            // -r 30 30帧/秒
            // scale 把原图修改下分辨率，缺少的地方不剪切不拉伸而是加黑边，再把所有处理后的图片二次处理成视频
            $shell = "ffmpeg -framerate 1/2 -start_number 1 -i \"$input_img_template\" -c:v libx264 -r 30 -vf \"scale=" . $max_width . ':' . $max_height . ':force_original_aspect_ratio=decrease,pad=' . $max_width . ':' . $max_height . ':(ow-iw)/2:(oh-ih)/2" -qscale 1 "' . $video_output . '" -y';
            $out = [];
            exec($shell, $out);
        }


        $music_output = "$video_output_path$key2-music.mp4";
        if (!file_exists($music_output)) {
            $audio_input_file = getRandomAudioFile();
            $shell = "ffmpeg -i \"$video_output\" -stream_loop -1 -i \"$audio_input_file\" -shortest -map 0:v:0 -map 1:a:0 -c:v copy \"$music_output\" -y";
            $out = [];
            exec($shell, $out);
        }


        foreach ($file_list as $filename) {
            rename($filename['new'], $filename['old']);
        }
    }
}
