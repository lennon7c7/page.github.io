#!/usr/bin/env php
<?php
/**
 * download bing img
 */
set_time_limit(0);


include './include/common.php';

$filepath = './images/bingimg/all/';
$url = 'https://www.bingimg.cn/list9999999999999999';
$haystack = get_curl($url);
if (empty($haystack)) {
    echo 'empty($haystack)' . PHP_EOL;
    return;
}


$matches = [];
preg_match_all('/<title>必应历史壁纸列表第([0-9]+)页 \| 必应高清壁纸 \| 必应每日美图<\/title>/', $haystack, $matches);
if (empty($matches[1][0])) {
    echo 'empty pagination' . PHP_EOL;
    return;
}


$max = $matches[1][0];
for ($i = 1; $i <= $max; $i++) {
    echo $i . PHP_EOL;
    $url = "https://www.bingimg.cn/list$i";
    $haystack = get_curl($url);
    if (empty($haystack)) {
        continue;
    }

    $matches = [];
    preg_match_all('/ src="(.*)" data-holder-rendered="true" class="card_img">/', $haystack, $matches);
    if (empty($matches[1])) {
        continue;
    }

    $imgs = $matches[1];
    foreach ($imgs as $old_img_url) {
        $new_img_url = str_replace('400x240', '1920x1080', $old_img_url);
        $new_img_url = str_replace('https://www.bingimg.cn/static/downimg/scale/SCALE.', 'https://www.bingimg.cn/down/', $new_img_url);

        $filename = str_replace('https://cn.bing.com/th?id=', '', $new_img_url);
        $filename = str_replace('https://www.bingimg.cn/down/', '', $filename);

        downloadImg($new_img_url, $filepath, $filename);
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
