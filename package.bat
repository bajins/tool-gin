1>1/* :::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
::  bajins 1.0.0  by bajins https://batch.bajins.com
:: 首发兼更新地址:https://batch.bajins.com
::
:: 使用时请将bajins.bat放入任意一个PATH中的目录以便调用
:: 但请确保bajins.bat拥有该目录的读写权限(因此最好不要选择system32)
:: 建议新建一个目录专供bajins.bat使用,再将这个目录添加到PATH中

@echo off
md "%~dp0$testAdmin$" 2>nul
if not exist "%~dp0$testAdmin$" (
    echo bajins不具备所在目录的写入权限! >&2
    exit /b 1
) else rd "%~dp0$testAdmin$"

setlocal enabledelayedexpansion

7za
if not "%errorlevel%" == "0" (
    :: cscript -nologo -e:jscript "%~f0" 这一段是执行命令，后面的是参数（组成方式：/key:value）
    :: %~f0 表示当前批处理的绝对路径,去掉引号的完整路径
    cscript -nologo -e:jscript "%~f0" https://woytu.github.io/files/7za.exe C:\Windows
)

set root=%~dp0
set files=%root%pyutils %root%static %root%templates

set project=key-gin

go get github.com/mitchellh/gox

gox


set otherList=_darwin_386,_darwin_amd64,_freebsd_386,_freebsd_amd64,_freebsd_arm,_netbsd_386,_netbsd_amd64,_netbsd_arm,_openbsd_386,_openbsd_amd64,_windows_386.exe,_windows_amd64.exe
:: 打包为zip
for %%i in (%otherList%) do (
    if exist "%root%%%i" (
        7za a %project%%%i.zip %files% %root%%%i
    )
)


set linuxList=_linux_386,_linux_amd64,_linux_arm,_linux_mips,_linux_mips64,_linux_mips64le,_linux_mipsle,_linux_s390x

:: 打包为tar.gz
for %%i in (%linuxList%) do (
    if exist "%root%%%i" (
        7za.exe a -ttar %project%%%i.tar %files% %root%%%i | 7za.exe a -tgzip %project%%%i.tar.gz %project%%%i.tar | del *.tar
    )
)


goto :EXIT

:EXIT
endlocal&exit /b %errorlevel%
*/

// ****************************  JavaScript  *******************************


var iRemote = WScript.Arguments(0);
iRemote = iRemote.toLowerCase();
var iLocal = WScript.Arguments(1);
iLocal = iLocal.toLowerCase()+"\\"+ iRemote.substring(iRemote.lastIndexOf("/") + 1);
var xPost = new ActiveXObject("Microsoft.XMLHTTP");
xPost.Open("GET", iRemote, 0);
xPost.Send();
var sGet = new ActiveXObject("ADODB.Stream");
sGet.Mode = 3;
sGet.Type = 1;
sGet.Open();
sGet.Write(xPost.responseBody);
sGet.SaveToFile(iLocal, 2);
sGet.Close();