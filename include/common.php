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
 * @param string $filepath
 * @param string $filename
 * @return void
 */
function downloadImg($url, $filepath, $filename)
{
    $filename_path = "$filepath$filename";

    if (!is_dir($filepath)) {
        $res = mkdir($filepath, 0777, true);
        if (!$res) {
            echo "mkdir $filepath error" . PHP_EOL;
            return;
        }
    }

    if (file_exists($filename_path) && filesize($filename_path) > 0) {
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

    file_put_contents($filename_path, $file_content);
    cleanExifInfo($filename_path);
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

            $temp_dir1[$file2] = $temp_dir2;
        }

        if (!empty($temp_dir1)) {
            $input_file[$file1] = $temp_dir1;
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
 * 重置图片文件大小
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
 * 删除顶部带水印的X轴部分
 * @param string $filename 文件名
 * @param int $watermark_px 水印高度
 * @return void
 */
function deleteTopWatermarkImage($filename, $watermark_px)
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

    imagejpeg($dst_image, __FUNCTION__ . "-$filename", 100);
}

/**
 * 删除底部带水印的X轴部分
 * @param string $filename 文件名
 * @param int $watermark_px 水印高度
 * @return void
 */
function deleteBottomWatermarkImage($filename, $watermark_px)
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

    imagejpeg($dst_image, __FUNCTION__ . "-$filename", 100);
}

/**
 * Convert image to jpeg image
 * @param string $old_filename 文件名
 * @return string
 */
function covertImage($old_filename)
{
    $filename_without_ext = pathinfo($old_filename, PATHINFO_FILENAME);
    $new_filename = "$filename_without_ext.jpg";

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

    $mix_width = 16;
    $mix_height = 9;

    $width = 160;
    $height = 90;
    foreach ($file_list as $filename) {
        if (empty($filename['width'])) {
            continue;
        }

        if ($filename['width'] > $width) {
            $width = $filename['width'];
            $height = $filename['height'];
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

    $mix_width = 9;
    $mix_height = 16;

    $width = 90;
    $height = 160;
    foreach ($file_list as $filename) {
        if (empty($filename['height'])) {
            continue;
        }

        if ($filename['height'] > $height) {
            $height = $filename['height'];
            $width = $filename['width'];
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
