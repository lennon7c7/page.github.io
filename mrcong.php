#!/usr/bin/env php
<?php
/**
 * download mrcong img
 */
set_time_limit(0);


include './include/simple_html_dom.php';
include './include/common.php';

$filepath_name = str_replace('.php', '', basename(__FILE__));
$filepath = "./images/$filepath_name/";
$url = 'https://mrcong.com';
$html_dom = file_get_html($url);
if (empty($html_dom)) {
    echo 'empty($html_dom)' . PHP_EOL;
    return;
}

foreach ($html_dom->find('.post-listing .post-box-title a') as $key => $tag_a) {
    if (empty($tag_a->href)) {
        continue;
    }

//    echo "$tag_a->innertext ";
    echo "$tag_a->href" . PHP_EOL;
    $html_dom2 = file_get_html($tag_a->href);
    if (empty($html_dom2)) {
        echo 'empty($html_dom2)' . PHP_EOL;
        continue;
    }

    $pre_index = 1;
    foreach ($html_dom2->find('a.shortc-button.medium.green') as $tag_a_2) {
        if (empty($tag_a_2->href)) {
            continue;
        }

        echo "  $tag_a_2->href" . PHP_EOL;
        $download_link = getMediafireDownloadLink($tag_a_2->href);
        if (!empty($download_link)) {
            // do something
            echo "    $download_link" . PHP_EOL;
        }
    }
}
