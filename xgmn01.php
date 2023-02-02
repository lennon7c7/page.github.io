#!/usr/bin/env php
<?php
/**
 * download xgmn01 img
 */
set_time_limit(0);


include './include/simple_html_dom.php';

$filepath_name = str_replace('.php', '', basename(__FILE__));
$filepath = "./images/$filepath_name/";
$pre = 'https://www.xgmn01.com';
$url = "$pre/new.html";
$html_dom = file_get_html($url);
if (empty($html_dom)) {
    echo 'empty($html_dom)' . PHP_EOL;
    return;
}

foreach ($html_dom->find('.widget-title a') as $element) {
    if (empty($element->href)) {
        continue;
    }

    echo $element->href . PHP_EOL;
    $html_dom2 = file_get_html($element->href);
    if (empty($html_dom2)) {
        echo 'empty($html_dom2)' . PHP_EOL;
        continue;
    }

    foreach ($html_dom2->find('.pagination a') as $element2) {
        if (empty($element2->href)) {
            continue;
        }

        echo "  $pre$element2->href" . PHP_EOL;
        $html_dom3 = file_get_html("$pre$element2->href");
        if (empty($html_dom3)) {
            echo 'empty($html_dom3)' . PHP_EOL;
            continue;
        }

        foreach ($html_dom3->find('.article-content img') as $element_img) {
            if (empty($element_img->src)) {
                continue;
            }

            $new_img_url = "$pre$element_img->src";
            echo "    $new_img_url" . PHP_EOL;

            $new_filepath = "$filepath$element_img->alt/";
            $new_filepath = str_replace('Xgyw.Net_', '', $new_filepath);

            $new_filename = basename($new_img_url);

            downloadImg($new_img_url, $new_filepath, $new_filename);
        }
    }
}


/**
 * @param string $url
 * @return string
 */
function get_curl($url)
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
            echo "目录 $filepath 创建失败" . PHP_EOL;
            return;
        }
    }

    if (file_exists($filename_path) && filesize($filename_path) > 0) {
        return;
    }

    file_put_contents($filename_path, file_get_contents($url));
}
