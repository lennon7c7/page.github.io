#!/usr/bin/env php
<?php
/**
 * download mrcong img
 */
set_time_limit(0);


include './include/simple_html_dom.php';
include './include/common.php';

$filepath_name = str_replace('.php', '', basename(__FILE__));
$GLOBALS['filepath'] = "./images/$filepath_name/";
$url = 'https://mrcong.com';
listPage($url);

function listPage($url)
{
    echo $url . PHP_EOL;
    $html_dom = file_get_html($url);
    if (empty($html_dom)) {
        echo 'empty($html_dom)' . PHP_EOL;
        return;
    }

    foreach ($html_dom->find('.post-listing .post-box-title a') as $tag_a) {
        if (!empty($tag_a->innertext) && is_dir("{$GLOBALS['filepath']}$tag_a->innertext")) {
            continue;
        }

        detailPage($tag_a->href);
    }

    $next_url = $html_dom->findFirst('head link[rel=next]')->href;
    if (!empty($next_url)) {
        listPage($next_url);
    }
}

function detailPage($url)
{
    if (empty($url)) {
        return;
    }

    echo "  $url" . PHP_EOL;
    $html_dom2 = file_get_html($url);
    if (empty($html_dom2)) {
        echo 'empty($html_dom2)' . PHP_EOL;
        return;
    }

    $title = $html_dom2->findFirst('#crumbs .current')->innertext;
    if (empty($title)) {
        $title = date('YmdHis');
    }
    $new_filepath = "{$GLOBALS['filepath']}$title/";

    foreach ($html_dom2->find('#fukie2 img.aligncenter') as $tag_img) {
        if (empty($tag_img->src)) {
            continue;
        }

        $new_img_url = "$tag_img->src";
        echo "    $new_img_url" . PHP_EOL;

        $new_filename = basename($new_img_url);
        $filename = "$new_filepath$new_filename";

        downloadImg($new_img_url, $filename);
    }


    // mediafire zip
//    foreach ($html_dom2->find('a.shortc-button.medium.green') as $tag_a_2) {
//        if (empty($tag_a_2->href)) {
//            continue;
//        }
//
//        echo "  $tag_a_2->href" . PHP_EOL;
//        $download_link = getMediafireDownloadLink($tag_a_2->href);
//        if (!empty($download_link)) {
//            // do something
//            echo "    $download_link" . PHP_EOL;
//        }
//    }


    $next_url = $html_dom2->findFirst('head link[rel=next]')->href;
    if (!empty($next_url)) {
        detailPage($next_url);
    }
}
