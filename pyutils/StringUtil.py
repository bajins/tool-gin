#!/usr/bin/env python
# -*- encoding: utf-8 -*-
# @Author : bajins www.bajins.com
# @File : StringUtil.py
# @Version: 1.0.0
# @Time : 2019/8/22 9:10
# @Project: windows-wallpaper-python
# @Package: 
# @Software: PyCharm
import json
import math
import random
import re
import string


def is_empty(obj):
    """
    判断数据是否为空
    :param obj:
    :return:
    """
    if isinstance(obj, str):
        if obj is None or len(obj) <= 0 or obj.strip() == '':
            return True

    elif isinstance(obj, set) or isinstance(obj, dict) or isinstance(obj, list):
        if obj is None or len(obj) <= 0 or bool(obj) or not any(obj):
            return True
    else:
        if obj or obj is None:
            return True
    return False


def not_empty(obj):
    """
    判断数据是否不为空
    :param obj:
    :return:
    """
    return not is_empty(obj)


def check_json(string):
    """
    确认是否为json
    :param string:
    :return:
    """
    try:
        json.loads(string)
        return True
    except BaseException as e:
        print(e)
        return False


def check_exist(string, substring):
    """
    判断字符串是否包含子串
    :param string:字符串
    :param substring: 子串
    :return:
    """
    string = string.lower()
    substring = substring.lower()
    # 使用正则表达式判断
    # if re.match("^.*" + substring + ".*", string):
    if string.find(substring) != -1 and substring in string:
        return True
    else:
        return False


def check_startswith(string, substring):
    """
    判断字符串是以什么开头
    :param string:字符串
    :param substring:需要判断的开头字符串
    :return:
    """
    # 检查你输入的是否是字符类型
    if isinstance(string, str):
        raise ValueError("参数不是字符串类型")
    # 判断字符串以什么开头
    if string.startswith(substring):
        return True

    return False


def check_endswith(string, substring):
    """
    判断字符串是以什么结尾
    :param string:字符串
    :param substring:需要判断的结尾字符串
    :return:
    """
    # 检查你输入的是否是字符类型
    if isinstance(string, str):
        raise ValueError("参数不是字符串类型")
    # 判断字符串以什么结尾
    if string.endswith(substring):
        return True

    return False


def random_lowercase_alphanumeric(length, charset="abcdefghijklmnopqrstuvwxyz0123456789_"):
    """
    生成一个指定长度的小写字母、数字、下划线的字符串
    :param length:
    :param charset:
    :return:
    """
    rd = ""
    for i in range(length):
        random_poz = math.floor(random.random() * len(charset))
        rd += charset[random_poz:random_poz + 1]
    return rd


def random_string(length=16):
    """
    生成一个指定长度的随机字符串，其中
    string.digits=0123456789
    string.ascii_letters=abcdefghigklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ
    """
    str_list = [random.choice(string.digits + string.ascii_letters) for i in range(length)]
    return ''.join(str_list)
