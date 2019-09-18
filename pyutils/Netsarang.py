#!/usr/bin/env python
# -*- encoding: utf-8 -*-
# @Author : bajins www.bajins.com
# @File : Netsarang.py
# @Version: 1.0.0
# @Time : 2019/9/17 18:22
# @Project: tool-gui-python
# @Package: 
# @Software: PyCharm
import base64
import sys
import time

from bs4 import BeautifulSoup

import Mail
import HttpUtil, StringUtil


def send_mail(mail, product):
    """
    根据产品和邮箱让Netsarang发送邮件
    :param mail:    邮箱
    :param product: 产品
    :return:
    """
    if product == "xshell":
        url = "https://www.netsarang.com/zh/xshell-download"

    if product == "xftp":
        url = "https://www.netsarang.com/zh/xftp-download"

    if product == "xmanager-power-suite":
        url = "https://www.netsarang.com/zh/xmanager-power-suite-download"

    if product == "xshell-plus":
        url = "https://www.netsarang.com/zh/xshell-plus-download"

    data = {"input[name='user-name']": mail.split("@")[0], "input[name='email']": mail}
    HttpUtil.crawling_selenium_bs_dictionary(url, data, 'input[value="开始试用"][type="submit"]')


def download(product):
    """
    获取下载链接地址
    :param product: 产品
    :return:
    """
    prefix = StringUtil.random_lowercase_alphanumeric(9)
    suffix = Mail.lin_shi_you_xiang_suffix()
    Mail.lin_shi_you_xiang_apply(prefix)
    mail = prefix + suffix
    send_mail(mail, product)
    time.sleep(5)
    mail_list = Mail.lin_shi_you_xiang_list(prefix)
    if len(mail_list) > 0:
        mail_content = Mail.lin_shi_you_xiang_get_mail(mail_list[0]["mailbox"], mail_list[0]["id"])
        html_text = base64.b64decode(mail_content.split("AmazonSES")[1])
        html = BeautifulSoup(html_text, features="html.parser")
        href = html.find("a", {"target": "_blank"}).text
        bs = HttpUtil.crawling_selenium(href)
        down_url = bs.find("a", {"target": "download_frame"})["href"]
        return down_url.replace(".exe", "r.exe")


if __name__ == '__main__':
    msg = download(sys.argv[1])
    print(msg)
