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
