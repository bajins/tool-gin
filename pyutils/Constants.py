#!/usr/bin/env python
# -*- encoding: utf-8 -*-
# @Author : bajins www.bajins.com
# @File : Constants.py
# @Version: 1.0.0
# @Time : 2019/8/21 15:32
# @Project: windows-wallpaper-python
# @Package:
# @Software: PyCharm

import os
import sys

APP_NAME = "Bajins工具"

# 解释器目录路径
EXECUTABLE_PATH = os.path.dirname(os.path.realpath(sys.executable))

# exe路径
APP_PATH = os.path.realpath(sys.argv[0])

# exe目录
APP_DIRECTORY = os.path.dirname(APP_PATH)

# 当前文件所在目录路径
CURRENT_PATH = os.path.dirname(__file__)

# 缓存目录路径
TEMP_PATH = sys.path[1]

# logo图片地址
LOGO_PATH = os.path.join(TEMP_PATH, "static", "logo.png")

# 配置路径
APP_CONF = os.path.join(APP_DIRECTORY, "app.conf")

HOSTS_PATH = "C:\\Windows\\System32\\drivers\\etc\\hosts"

CHINAZ_DNS = "http://tool.chinaz.com/dns"

MYSSL_DNS = "https://myssl.com/api/v1/tools/dns_query"

SHORT_TIME_MAIL = "https://shorttimemail.com"

SHORT_TIME_MAIL_DNS = SHORT_TIME_MAIL + "/net/dns/query"

LIN_SHI_YOU_XIANG = "https://www.linshiyouxiang.net"

GITHUB_DOMAIN = [
    "assets-cdn.github.com",
    "avatars.githubusercontent.com",
    "avatars0.githubusercontent.com",
    "avatars1.githubusercontent.com",
    "codeload.github.com",
    "documentcloud.github.com",
    "gist.github.com",
    "github.com",
    "github.global.ssl.fastly.net",
    "github.io",
    "github-cloud.s3.amazonaws.com",
    "global-ssl.fastly.net",
    "help.github.com",
    "nodeload.github.com",
    "raw.github.com",
    "status.github.com",
    "training.github.com",
    "www.github.com"
]
