#!/usr/bin/env php
<?php
/**
 * download xgmn01 img
 */
set_time_limit(0);


include './include/simple_html_dom.php';
include './include/common.php';

$filepath_name = str_replace('.php', '', basename(__FILE__));
$filepath = "./images/$filepath_name/";
$pre = 'https://www.xgmn01.com';
$url = "$pre/new.html";
$html_dom = file_get_html($url);
if (empty($html_dom)) {
    echo 'empty($html_dom)' . PHP_EOL;
    return;
}

foreach ($html_dom->find('.widget-title a') as $key => $tag_a) {
    if (empty($tag_a->href)) {
        continue;
    }

    echo $tag_a->href . PHP_EOL;
    if (!empty($tag_a->title) && is_dir("$filepath$tag_a->title")) {
        continue;
    }

    $html_dom2 = file_get_html($tag_a->href);
    if (empty($html_dom2)) {
        echo 'empty($html_dom2)' . PHP_EOL;
        continue;
    }

    $detail_urls = [];
    foreach ($html_dom2->find('.pagination a') as $tag_a_2) {
        if (empty($tag_a_2->href)) {
            continue;
        }

        $detail_urls[] = "$pre$tag_a_2->href";
    }
    $detail_urls = array_unique($detail_urls);

    $pre_index = 1;
    foreach ($detail_urls as $detail_url) {
        if (empty($detail_url)) {
            continue;
        }

        echo "  $detail_url" . PHP_EOL;
        $html_dom3 = file_get_html($detail_url);
        if (empty($html_dom3)) {
            echo 'empty($html_dom3)' . PHP_EOL;
            continue;
        }

        foreach ($html_dom3->find('.article-content img') as $tag_img) {
            if (empty($tag_img->src)) {
                continue;
            }

            $new_img_url = "$pre$tag_img->src";
            echo "    $new_img_url" . PHP_EOL;

            $new_filepath = "$filepath$tag_a->title/";

            $new_filename = $pre_index . '-' . basename($new_img_url);
            $pre_index++;

            downloadImg($new_img_url, $new_filepath, $new_filename);
        }
    }
}
