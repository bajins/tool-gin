#!/usr/bin/env python
# -*- encoding: utf-8 -*-
# @Author : bajins www.bajins.com
# @File : HttpUtil.py
# @Version: 1.0.0
# @Time : 2019/8/21 15:32
# @Project: windows-wallpaper-python
# @Package:
# @Software: PyCharm


import json
import os
import socket
import sys
import time
import urllib

import requests
import urllib3
from bs4 import BeautifulSoup
from selenium import webdriver
from selenium.webdriver.common.by import By

import FileUtil


def get(url, data):
    """
    get请求
    :param url:请求地址
    :param data:数据，map或dict格式
    :return:
    """
    return requests.get(url=url, params=data)


def post(url, data):
    """
    get请求
    :param url:请求地址
    :param data:数据，map或dict格式
    :return:
    """
    return requests.post(url=url, params=data)


def delete(url, data):
    """
    delete请求
    :param url:请求地址
    :param data:数据，map或dict格式
    :return:
    """
    return requests.delete(url=url, params=data)


def get_json(url, data):
    """
    get请求返回结果转json
    :param url:
    :param data:
    :return:
    """
    return json.loads(requests.get(url=url, params=data).text)


def crawling_bs(url, data):
    """
    https://www.crummy.com/software/BeautifulSoup/bs4/doc.zh/#id9
    使用BeautifulSoup库爬取数据，
    解析器有：html.parser、lxml、xml、html5lib
    推荐使用lxml作为解析器，速度快，容错能力强，效率高
    :param url:
    :param data:
    :return:
    """
    resp = requests.get(url=url, params=data)
    return BeautifulSoup(resp.text, features="lxml")


def download_chromedriver():
    """
    下载chrome驱动
    http://chromedriver.storage.googleapis.com/index.html
    :return:
    """
    # 获取版本号列表
    url = "http://chromedriver.storage.googleapis.com/"
    result = crawling_bs(url, {"delimiter": "/", "prefix": ""})
    prefix = result.find_all("prefix")
    # 过滤
    # info = [s.extract() for s in prefix('prefix')]
    ver = []
    for s in prefix:
        t = s.text
        # 判断如果全是字母就不是版本号
        if not t.replace("/", "").isalpha():
            ver.append(t)
    # 对版本号降序排序
    ver.sort(reverse=True)

    # 获取版本下面的文件列表
    driver_list = crawling_bs(url, {"delimiter": "/", "prefix": ver[0]})
    filename_list = driver_list.find_all("key")

    for s in filename_list:
        s = s.text
        # 如果在文件名中找到系统平台名称
        if s.find(sys.platform) != -1:
            filename = s[len(ver[0]):]
            # 下载文件
            download_file(url + s, None, filename)
            FileUtil.zip_extract(filename, None)


def download_taobao_chromedriver():
    """
    下载淘宝镜像chromedriver
    http://npm.taobao.org/mirrors/chromedriver
    :return:
    """
    # 获取版本号列表
    url = "http://npm.taobao.org/mirrors/chromedriver/"
    result = crawling_bs(url, None)
    prefix = result.find("pre").find_all("a")
    # 过滤
    # info = [s.extract() for s in prefix('prefix')]
    ver = []
    for s in prefix:
        t = s.text
        # 判断如果全是字母就不是版本号
        if not t.replace("/", "").isalpha() and t.endswith("/"):
            ver.append(t)
    # 对版本号降序排序
    ver.sort(reverse=True)

    latestVersionUrl = url + ver[0]
    # 获取版本下面的文件列表
    driver_list = crawling_bs(latestVersionUrl, None)
    filename_list = driver_list.find("pre").find_all("a")

    for s in filename_list:
        s = s.text
        # 如果在文件名中找到系统平台名称
        if s.find(sys.platform) != -1:
            # 下载文件
            download_file(latestVersionUrl + s, None, s)
            FileUtil.zip_extract(s, None)


def selenium_driver(url):
    """
    获取驱动
    :param url:
    :return:
    """
    if sys.platform == "win32":
        path = "./chromedriver.exe"
    else:
        path = "./chromedriver"

    if not os.path.exists(path):
        download_taobao_chromedriver()

    # chrome选项
    options = webdriver.ChromeOptions()
    # 设置chrome浏览器无界面模式
    options.add_argument('--headless')
    # 解决DevToolsActivePort文件不存在的报错
    options.add_argument('--no-sandbox')
    # 指定浏览器分辨率
    options.add_argument('window-size=1600x900')
    # 谷歌文档提到需要加上这个属性来规避bug
    options.add_argument('--disable-gpu')
    # 隐藏滚动条, 应对一些特殊页面
    options.add_argument('--hide-scrollbars')
    # 不加载图片, 提升速度
    options.add_argument('blink-settings=imagesEnabled=false')
    # 打开浏览器,executable_path指定驱动位置
    driver = webdriver.Chrome(chrome_options=options, executable_path=path)
    # 最大化浏览器
    # driver.maximize_window()
    # 最小化浏览器
    driver.minimize_window()
    # 打开网站
    driver.get(url)
    return driver


def crawling_selenium(url):
    """
    使用selenium库打开一个链接并获取网页源码，
    再利用BeautifulSoup操作数据
    :param url:
    :return:
    """
    try:
        driver = selenium_driver(url)
        # 获取网页源代码
        html = driver.page_source
        # 使用BeautifulSoup创建html代码的BeautifulSoup实例
        return BeautifulSoup(html, features="html.parser")
    finally:
        # 关闭当前窗口。
        driver.close()
        # 关闭浏览器并关闭chreomedriver进程
        driver.quit()


def crawling_selenium_bs(url, input_el, input_text):
    """
    使用selenium库打开一个链接并获取网页源码，
    再利用BeautifulSoup操作数据
    :param url:
    :param input_el: input标签的id，name或class
    :param input_text: 输入内容
    :return:
    """
    try:
        driver = selenium_driver(url)
        # n次点击加载更多
        # for i in range(0, 5):
        #     # 点击加载更多
        #     driver.find_element_by_class_name("home-news-footer").click()
        #     # 找到加载更多按钮，点击
        #     driver.find_element(By.LINK_TEXT, "加载更多").click()
        #     # 延时两秒
        #     time.sleep(2)
        # 使用selenium通过id，name或class的方式来获取到这个input标签
        input_element = driver.find_element_by_class_name(input_el)
        # 传入值，输入的内容
        input_element.send_keys(input_text)
        # 提交
        input_element.submit()
        # 延时
        time.sleep(2)
        # 获取网页源代码
        html = driver.page_source
        # 使用BeautifulSoup创建html代码的BeautifulSoup实例
        return BeautifulSoup(html, features="html.parser")
    finally:
        # 关闭当前窗口。
        driver.close()
        # 关闭浏览器并关闭chreomedriver进程
        driver.quit()


def crawling_selenium_bs_dictionary(url, input_dictionary, click_btn):
    """
    使用selenium库打开一个链接并获取网页源码，
    再利用BeautifulSoup操作数据
    :param url:                 访问链接
    :param input_dictionary:    input标签和内容
    :param click_btn:           点击提交的按钮
    :return:
    """
    try:
        driver = selenium_driver(url)
        # 使用selenium通过id，name或class的方式来获取到这个input标签
        for key, value in input_dictionary.items():
            # 查找元素，传入值（输入的内容）
            driver.find_element_by_css_selector(key).send_keys(value)

        # 提交
        # driver.find_element_by_xpath(click_btn).click()
        driver.find_element_by_css_selector(click_btn).click()
        # driver.find_element_by_class_name(click_btn).click()
        # 延时
        time.sleep(2)
        # 获取网页源代码
        html = driver.page_source
        # 使用BeautifulSoup创建html代码的BeautifulSoup实例
        return BeautifulSoup(html, features="html.parser")
    finally:
        # 关闭当前窗口。
        driver.close()
        # 关闭浏览器并关闭chreomedriver进程
        driver.quit()


def download_file(url, mkdir, name=""):
    """
    用requests下载文件
    :param url:
    :param mkdir:
    :param name:
    :return:
    """
    # detectionModule("requests")
    # 判断文件名称是否传入
    if name is None or name == "":
        ur = str(url).split("/")
        # 如果没传，就取URL中最后的文件名
        name = ur[len(ur) - 1]

    # 判断是否传入文件夹
    if mkdir is not None and mkdir != "":
        # 判断目录是否存在
        if not os.path.exists(mkdir):
            # 目录不存在则创建
            os.mkdir(mkdir)
        name = os.path.join(mkdir, name)

    # 判断文件是否存在
    # if not os.path.exists(name):
    if not os.path.isfile(name):
        # 去除警告
        requests.packages.urllib3.disable_warnings()
        requests.adapters.DEFAULT_RETRIES = 5
        # 打开session
        s = requests.session()
        s.keep_alive = False
        headers = {
            "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) "
                          "Chrome/72.0.3626.109 Safari/537.36 "
        }
        # 文件不存在才保存
        with open(name, "wb") as f:
            f.write(s.get(url, headers=headers, verify=False, timeout=30).content)


def download_file_list(urls, mkdir, name):
    """
    用urllib批量下载文件
    :param urls:
    :param mkdir:
    :param name:
    :return:
    """
    # 老版本去除警告方法
    # from requests.packages.urllib3.exceptions import InsecureRequestWarning
    # requests.packages.disable_warnings(InsecureRequestWarning)

    # 新版去除警告方法
    urllib3.disable_warnings(urllib3.exceptions.InsecureRequestWarning)
    for url in urls:
        # 判断文件名称是否传入
        if name.strip() == '':
            ur = str(url).split("/")
            # 如果没传，就取URL中最后的文件名
            name = ur[len(ur) - 1]
        # 判断是否传入文件夹
        if mkdir.strip() != '':
            # 判断目录是否存在
            if not os.path.exists(mkdir):
                # 目录不存在则创建
                os.mkdir(mkdir)
            name = mkdir + name
        # os.path.join将多个路径组合后返回
        # LocalPath = os.path.join('C:/Users/goatbishop/Desktop',file)
        # 第一个参数url:需要下载的网络资源的URL地址
        # 第二个参数LocalPath:文件下载到本地后的路径
        urllib.request.urlretrieve(url, name)
        # response = urllib.request.urlopen(url)
        # pic = response.read()
        # with open(name, 'wb') as f:
        #     f.write(pic)


def get_host_ip():
    """
    查询本机ip地址
    :return: ip
    """
    try:
        s = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
        s.connect(('8.8.8.8', 80))
        ip = s.getsockname()[0]
    finally:
        s.close()

    return ip


def get_remote_ip(host_name):
    """
    获取指定域名IP地址
    :param host_name:域名
    :return:
    """
    try:
        return socket.gethostbyname(host_name)
    except BaseException as e:
        print(" %s:%s" % (host_name, e))


if __name__ == '__main__':
    download_taobao_chromedriver()
