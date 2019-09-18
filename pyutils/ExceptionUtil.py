#!/usr/bin/env python
# -*- encoding: utf-8 -*-
# @Author : bajins www.bajins.com
# @File : ExceptionUtil.py
# @Version: 1.0.0
# @Time : 2019/8/21 15:32
# @Project: windows-wallpaper-python
# @Package:
# @Software: PyCharm


class MsgException(BaseException):
    """
      自定义异常
      使用 raise MsgException('请求接口错误')
    """

    def __init__(self, msg):
        self.msg = msg

    def __str__(self):
        return self.msg
