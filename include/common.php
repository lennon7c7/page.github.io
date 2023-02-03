<?php

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

    file_put_contents($filename_path, file_get_contents($url));

    cleanExifInfo($filename_path);
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
