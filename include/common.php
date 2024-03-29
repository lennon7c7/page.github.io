<?php

/**
 * @param string $url
 * @return string
 */
function curl_get($url)
{
    $curl = curl_init();

    curl_setopt_array($curl, [
        CURLOPT_URL => $url,
        CURLOPT_RETURNTRANSFER => true,
        CURLOPT_ENCODING => '',
        CURLOPT_MAXREDIRS => 10,
        CURLOPT_TIMEOUT => 0,
        CURLOPT_FOLLOWLOCATION => true,
        CURLOPT_HTTP_VERSION => CURL_HTTP_VERSION_1_1,
        CURLOPT_CUSTOMREQUEST => 'GET',

        CURLOPT_SSL_VERIFYPEER => false,
        CURLOPT_SSL_VERIFYHOST => false,
        CURLOPT_USERAGENT => 'Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/45.0.2454.98 Safari/537.36',
    ]);

    $response = curl_exec($curl);

    curl_close($curl);

    return $response;
}

/**
 * @param string $url
 * @param string $filename
 * @return void
 */
function downloadImg($url, $filename)
{
    $filepath = dirname($filename);

    if (!is_dir($filepath)) {
        $res = mkdir($filepath, 0777, true);
        if (!$res) {
            echo "mkdir $filepath error" . PHP_EOL;
            return;
        }
    }

    if (file_exists($filename) && filesize($filename) > 0) {
        return;
    }

    $opts = [
        'http' => [
            'method' => 'GET',
            'timeout' => 5,
        ]
    ];
    $count = 0;
    while ($count < 3 && ($file_content = @file_get_contents($url, false, stream_context_create($opts))) === FALSE) $count++;
    if (empty($file_content)) {
        return;
    }

    file_put_contents($filename, $file_content);
    $filename = covertImage($filename);
    resizeImageToEvenNumber($filename);
    deleteTopWatermarkImage($filename);
    deleteBottomWatermarkImage($filename);
//    cleanExifInfo($filename);
}

/**
 * @param string $page_link
 * @return string download link
 */
function getMediafireDownloadLink($page_link)
{
    $response = curl_get($page_link);

    $matches = [];
    preg_match_all('/{ window.location.href = \'(.*)\'; }/', $response, $matches);
    if (!empty($matches[1][0])) {
        $download_link = $matches[1][0];

        $header_array = get_headers($download_link, true);
        $size = bcdiv($header_array['content-length'], 1073741824, 2);
        $download_link .= "?size={$size}G";

        return $download_link;
    }

    $html_dom = str_get_html($response);
    if (!empty($html_dom->getElementById('downloadButton')->href)) {
        $download_link = $html_dom->getElementById('downloadButton')->href;

        $header_array = get_headers($download_link, true);
        $size = bcdiv($header_array['content-length'], 1073741824, 2);
        $download_link .= "?size={$size}G";

        return $download_link;
    }

    return '';
}

/**
 * @param string $path
 * @return string abs file path
 */
function getRandomAudioFile()
{
    $path = './audio/';
    $keep_needle = ['m4a', 'flac', 'mp3', 'wav', 'wma', 'aac'];
    $files = [];
    foreach (scandir($path) as $file) {
        $filename_ext = pathinfo($file, PATHINFO_EXTENSION);
        if (!in_array($filename_ext, $keep_needle)) {
            continue;
        }

        $files[] = "$path$file";
    }

    return $files[array_rand($files)];
}

/**
 * @return array input struction data
 */
function getInputImg()
{
    $top_dir = './images';
    $filter_file = ['.', '..', 'desktop.ini'];
    $input_file = [];
    foreach (scandir($top_dir) as $file1) {
        $dir1 = "$top_dir/$file1";
        if (in_array($file1, $filter_file) || !is_dir($dir1)) {
            continue;
        }

        $temp_dir1 = [];
        foreach (scandir($dir1) as $file2) {
            $dir2 = "$top_dir/$file1/$file2";
            if (in_array($file2, $filter_file) || !is_dir($dir2)) {
                continue;
            }

            $temp_dir2 = [];
            foreach (scandir($dir2) as $file3) {
                $dir3 = "$top_dir/$file1/$file2/$file3";
                if (in_array($file3, $filter_file) || !is_dir($dir2) || !isImgFileExt($dir3)) {
                    continue;
                }

                $temp_dir2[] = $dir3;
            }

            if (!empty($temp_dir2)) {
                $temp_dir1[$dir2] = $temp_dir2;
            }
        }

        if (!empty($temp_dir1)) {
            $input_file[$dir1] = $temp_dir1;
        }
    }

    return $input_file;
}

/**
 * @param string $filename
 * @return bool
 */
function isImgFileExt($filename)
{
    $keep_needle = ['jpg', 'jpeg', 'gif', 'png', 'bmp', 'webp'];
    foreach ($keep_needle as $ext) {
        $filename_ext = pathinfo($filename, PATHINFO_EXTENSION);
        if ($filename_ext === $ext) {
            return true;
        }
    }

    return false;
}

/**
 * 重置图片尺寸
 * @param string $filename 文件名
 * @param string $dst_width 修改后最大宽度
 * @param string $dst_height 修改后最大高度
 * @return void
 */
function resizeImage($filename, $dst_width, $dst_height)
{
    $ext = explode('.', $filename);
    $ext = $ext[count($ext) - 1];

    if ($ext == 'jpg' || $ext == 'jpeg') {
        $src_image = imagecreatefromjpeg($filename);
    } elseif ($ext == 'gif') {
        $src_image = imagecreatefromgif($filename);
    } elseif ($ext == 'png') {
        $src_image = imagecreatefrompng($filename);
    } elseif ($ext == 'bmp') {
        $src_image = imagecreatefrombmp($filename);
    } elseif ($ext == 'webp') {
        $src_image = imagecreatefromwebp($filename);
    }

    if (empty($src_image)) {
        return;
    }

    $src_width = imagesx($src_image);
    $src_height = imagesy($src_image);

    $dst_image = imagecreatetruecolor($dst_width, $dst_height);
    imagecopyresized($dst_image, $src_image, 0, 0, 0, 0, $dst_width, $dst_height, $src_width, $src_height);

    imagejpeg($dst_image, $filename, 100);
}

/**
 * 修改图片尺寸成能够被2所整除的整数
 * @param string $filename 文件名
 * @return void
 */
function resizeImageToEvenNumber($filename)
{
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
}

/**
 * 序列化某目录下的jpg文件
 * @param string $dir
 * @return void
 */
function serializeJpgFilename($dir)
{
    $temp = date('/temp-YmdHis-');
    $keep_needle = ['jpg'];
    $files = [];
    foreach (scandir($dir) as $file) {
        $filename_ext = pathinfo($file, PATHINFO_EXTENSION);
        if (!in_array($filename_ext, $keep_needle)) {
            continue;
        }

        $files[] = "$dir/$file";
    }
    foreach ($files as $key => $file) {
        $new_file = dirname($file) . $temp . sprintf('%04d', $key + 1) . '.jpg';
        rename($file, $new_file);
    }

    $files = [];
    foreach (scandir($dir) as $file) {
        $filename_ext = pathinfo($file, PATHINFO_EXTENSION);
        if (!in_array($filename_ext, $keep_needle)) {
            continue;
        }

        $files[] = "$dir/$file";
    }
    foreach ($files as $key => $file) {
        $new_file = dirname($file) . '/' . sprintf('%04d', $key + 1) . '.jpg';
        rename($file, $new_file);
    }
}

/**
 * 图像 - 文字识别
 * @param string $filename
 * @return void
 */
function ocr_xgmn01($filename)
{
    if (!file_exists($filename)) {
        return;
    }

    list($src_w, $src_h) = getimagesize($filename);

    $zoom = 3;
    $dst_w = bcdiv($src_w, $zoom);
    $dst_h = bcdiv($src_h, $zoom);

    $dst_scale = $dst_h / $dst_w; // 目标图像长宽比
    $src_scale = $src_h / $src_w; // 原图长宽比

    if ($src_scale >= $dst_scale) {  // 过高
        $w = intval($src_w);
        $h = intval($dst_scale * $w);

        $x = 0;
        $y = ($src_h - $h) / 3;
    } else { // 过宽
        $h = intval($src_h);
        $w = intval($h / $dst_scale);

        $x = ($src_w - $w) / 2;
        $y = 0;
    }

    // 剪裁
    $source = imagecreatefromjpeg($filename);
    $croped = imagecreatetruecolor($w, $h);
    imagecopy($croped, $source, 0, 0, $x, $y, $src_w, $src_h);

    // 缩放
    $scale = $dst_w / $w;
    $target = imagecreatetruecolor($dst_w, $dst_h);
    $final_w = intval($w * $scale);
    $final_h = intval($h * $scale);
    imagecopyresampled($target, $croped, 0, 0, 0, 0, $final_w, $final_h, $w, $h);

    // 保存
    $timestamp = time();
    $temp_file = "$timestamp.jpg";
    imagejpeg($target, $temp_file);
    imagedestroy($target);

    $img = file_get_contents($temp_file);
    $base64 = base64_encode($img);
    unlink($temp_file);

    $post_field = ['image' => $base64];
    $curl = curl_init();

    curl_setopt_array($curl, [
        CURLOPT_URL => 'https://www.paddlepaddle.org.cn/paddlehub-api/image_classification/chinese_ocr_db_crnn_mobile',
        CURLOPT_RETURNTRANSFER => true,
        CURLOPT_ENCODING => '',
        CURLOPT_MAXREDIRS => 10,
        CURLOPT_TIMEOUT => 10,
        CURLOPT_FOLLOWLOCATION => true,
        CURLOPT_HTTP_VERSION => CURL_HTTP_VERSION_1_1,
        CURLOPT_CUSTOMREQUEST => 'POST',
        CURLOPT_POSTFIELDS => json_encode($post_field),

        CURLOPT_SSL_VERIFYPEER => false,
        CURLOPT_SSL_VERIFYHOST => false,
        CURLOPT_USERAGENT => 'Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/45.0.2454.98 Safari/537.36',
    ]);

    $response = curl_exec($curl);
    curl_close($curl);

    if (empty($response)) {
        return;
    }

    $response_array = json_decode($response, true);
    if (empty($response_array['result'][0]['data'])) {
        return;
    }

    $matchWatermark = ['www.xgyw.net'];
    foreach ($response_array['result'][0]['data'] as $value) {
        if (in_array(strtolower($value['text']), $matchWatermark)) {
            unlink($filename);
            return;
        }
    }
}


/**
 * 删除顶部带水印的X轴部分
 * @param string $filename 文件名
 * @param int $watermark_px 水印高度
 * @return void
 */
function deleteTopWatermarkImage($filename, $watermark_px = 100)
{
    $ext = explode('.', $filename);
    $ext = $ext[count($ext) - 1];

    if ($ext == 'jpg' || $ext == 'jpeg') {
        $src_image = imagecreatefromjpeg($filename);
    } elseif ($ext == 'gif') {
        $src_image = imagecreatefromgif($filename);
    } elseif ($ext == 'png') {
        $src_image = imagecreatefrompng($filename);
    } elseif ($ext == 'bmp') {
        $src_image = imagecreatefrombmp($filename);
    } elseif ($ext == 'webp') {
        $src_image = imagecreatefromwebp($filename);
    }

    if (empty($src_image)) {
        return;
    }

    $src_width = imagesx($src_image);
    $src_height = imagesy($src_image);

    $dst_width = $src_width - $watermark_px;
    $dst_height = $src_height - $watermark_px;

    $dst_image = imagecreatetruecolor($dst_width, $dst_height);
    imagecopyresized($dst_image, $src_image, 0, 0, $watermark_px, $watermark_px, $src_width, $src_height, $src_width, $src_height);

    imagejpeg($dst_image, $filename, 100);
}

/**
 * 删除底部带水印的X轴部分
 * @param string $filename 文件名
 * @param int $watermark_px 水印高度
 * @return void
 */
function deleteBottomWatermarkImage($filename, $watermark_px = 100)
{
    $ext = explode('.', $filename);
    $ext = $ext[count($ext) - 1];

    if ($ext == 'jpg' || $ext == 'jpeg') {
        $src_image = imagecreatefromjpeg($filename);
    } elseif ($ext == 'gif') {
        $src_image = imagecreatefromgif($filename);
    } elseif ($ext == 'png') {
        $src_image = imagecreatefrompng($filename);
    } elseif ($ext == 'bmp') {
        $src_image = imagecreatefrombmp($filename);
    } elseif ($ext == 'webp') {
        $src_image = imagecreatefromwebp($filename);
    }

    if (empty($src_image)) {
        return;
    }

    $src_width = imagesx($src_image);
    $src_height = imagesy($src_image);

    $dst_width = $src_width - $watermark_px;
    $dst_height = $src_height - $watermark_px;

    $dst_image = imagecreatetruecolor($dst_width, $dst_height);
    imagecopyresized($dst_image, $src_image, 0, 0, 0, 0, $src_width, $src_height, $src_width, $src_height);

    imagejpeg($dst_image, $filename, 100);
}

/**
 * Convert image to jpeg image
 * @param string $old_filename 文件名
 * @return string
 */
function covertImage($old_filename)
{
    $filename_without_ext = pathinfo($old_filename, PATHINFO_FILENAME);
    $new_filename = dirname($old_filename) . "/$filename_without_ext.jpg";

    $ext = explode('.', $old_filename);
    $ext = $ext[count($ext) - 1];

    if ($ext == 'jpg') {
        return $old_filename;
    } elseif ($ext == 'jpeg') {
        rename($old_filename, $new_filename);
        return $new_filename;
    } elseif ($ext == 'gif') {
        $src_image = imagecreatefromgif($old_filename);
    } elseif ($ext == 'png') {
        $src_image = imagecreatefrompng($old_filename);
    } elseif ($ext == 'bmp') {
        $src_image = imagecreatefrombmp($old_filename);
    } elseif ($ext == 'webp') {
        $src_image = imagecreatefromwebp($old_filename);
    }

    if (empty($src_image)) {
        return '';
    }

    imagejpeg($src_image, $new_filename, 100);
    unlink($old_filename);

    return $new_filename;
}

/**
 * 获取当前文件列表
 * @return array
 */
function getCurrentFileList()
{
    $files = [];

    // 以当前目录的方式转
    $shell = "dir /b";
    $output = [];
    exec($shell, $output);

    foreach ($output as $filename) {
        $needle = '.';
        if (!stristr($filename, $needle)) {
            continue;
        }

        $needle =
            /** @lang text */
            '<DIR>';
        if (stristr($filename, $needle)) {
            continue;
        }

        $files[] = $filename;
    }

    return $files;
}

/**
 * 获取适用PC的长方形尺寸
 * @param array $file_list
 * @return array
 */
function getPCRectangle($file_list)
{
    $max_width = 2560;
    $max_height = 1440;

    $mix_width = 160;
    $mix_height = 90;

    $width = 160;
    $height = 90;
    foreach ($file_list as $filename) {
        $src_image = imagecreatefromjpeg($filename);
        $src_width = imagesx($src_image);
        $src_height = imagesy($src_image);

        if (empty($src_width)) {
            continue;
        }

        if ($src_width > $width) {
            $width = $src_width;
            $height = $src_height;
        }
    }

    if ($width > $max_width || $height > $max_height) {
        return [$max_width, $max_height];
    } elseif ($width < $mix_width || $height < $mix_height) {
        return [$mix_width, $mix_height];
    } else {
        $height = bcmul(bcdiv($width, $mix_width), $mix_height);
        return [$width, $height];
    }
}

/**
 * 获取适用phone的长方形尺寸
 * @param array $file_list
 * @return array
 */
function getPhoneRectangle($file_list)
{
    $max_width = 1440;
    $max_height = 2560;

    $mix_width = 90;
    $mix_height = 160;

    $width = 90;
    $height = 160;
    foreach ($file_list as $filename) {
        $src_image = imagecreatefromjpeg($filename);
        $src_width = imagesx($src_image);
        $src_height = imagesy($src_image);

        if (empty($src_height)) {
            continue;
        }

        if ($src_height > $height) {
            $height = $src_height;
            $width = $src_width;
        }
    }

    if ($width > $max_width || $height > $max_height) {
        return [$max_width, $max_height];
    } elseif ($width < $mix_width || $height < $mix_height) {
        return [$mix_width, $mix_height];
    } else {
        $width = bcmul(bcdiv($height, $mix_height), $mix_width);
        return [$width, $height];
    }
}

/**
 * @param string $filename
 */
function cleanExifInfo($filename)
{
    $shell = 'exiftool -All= -overwrite_original -m -q -q -Title="" -Description="" -Subject="" -Creator="" -LastKeywordXMP="" ' . $filename;
    $out = [];
    exec($shell, $out);

    if (!empty($out)) {
        print_r($out);
    }
}

/**
 * 统一图片尺寸
 * @param string $dir 图片目录
 * @return void
 */
function uniformImageSizeByDir($dir)
{
    /**
     * 统一图片尺寸
     * @param int $canvas_width 画布宽
     * @param int $canvas_height 画布高
     * @param string $watermark_filename 水印GDImage
     * @return void
     */
    function uniform_image_size($canvas_width, $canvas_height, $watermark_filename)
    {
        // 获取水印图片的宽高
        list($watermark_width, $watermark_height) = getimagesize($watermark_filename);
        if ($watermark_width == $canvas_width && $watermark_height == $canvas_width) {
            return;
        }

        $canvas = imagecreatetruecolor($canvas_width, $canvas_height);

        // 创建图片的实例
        $watermark_img = imagecreatefromstring(file_get_contents($watermark_filename));

        if ($canvas_width == $watermark_width) {
            // 将水印图片水平居中
            imagecopy($canvas, $watermark_img, 0, $canvas_height / 2 - $watermark_height / 2, 0, 0, $watermark_width, $watermark_height);
        } elseif ($canvas_height == $watermark_height) {
            // 将水印图片垂直居中
            imagecopy($canvas, $watermark_img, $canvas_height / 2 - $watermark_height / 2, 0, 0, 0, $watermark_width, $watermark_height);
        } else {
            // 将水印图片水平、垂直居中
            imagecopy($canvas, $watermark_img, $canvas_height / 2 - $watermark_height / 2, $canvas_height / 2 - $watermark_height / 2, 0, 0, $watermark_width, $watermark_height);
        }

        imagejpeg($canvas, $watermark_filename);
        imagedestroy($watermark_img);
        imagedestroy($canvas);
    }

    $keep_needle = ['jpg'];
    $files = [];
    $max_width = 160;
    $max_height = 90;
    foreach (scandir($dir) as $file) {
        $filename_ext = pathinfo($file, PATHINFO_EXTENSION);
        if (!in_array($filename_ext, $keep_needle)) {
            continue;
        }

        $old_filename = "$dir/$file";
        list($canvas_width, $canvas_height) = getimagesize($old_filename);
        $max_width = max($max_width, $canvas_width);
        $max_height = max($max_height, $canvas_height);
        $files[] = $old_filename;
    }

    foreach ($files as $file) {
        uniform_image_size($max_width, $max_height, $file);
    }
}
