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

foreach ($html_dom->find('.widget-title a') as $key => $element) {
    if (empty($element->href)) {
        continue;
    }

    echo $element->href . PHP_EOL;
    $html_dom2 = file_get_html($element->href);
    if (empty($html_dom2)) {
        echo 'empty($html_dom2)' . PHP_EOL;
        continue;
    }

    $pre_index = 1;
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

            $new_filename = $pre_index . '-' . basename($new_img_url);
            $pre_index++;

            downloadImg($new_img_url, $new_filepath, $new_filename);
        }
    }

    // todo dev fast
//    break;
//    if ($key > 0) {
//        break;
//    }
}
