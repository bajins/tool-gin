/**
 *
 * @Description:
 * @Author: https://www.bajins.com
 * @File: reptile_test.go
 * @Version: 1.0.0
 * @Time: 2019/9/19 11:13
 * @Project: tool-gin
 * @Package:
 * @Software: GoLand
 */
package reptile

import (
	"bytes"
	"fmt"
	"github.com/chromedp/cdproto/target"
	"github.com/chromedp/chromedp"
	"log"
	"net/mail"
	"strings"
	"testing"
	"time"
)

func TestLinShiYouXiangSuffix(t *testing.T) {
	LinShiYouXiangSuffix()
}

func TestLinShiYouXiangList(t *testing.T) {
	list, _ := LinShiYouXiangList("5wij52emu")
	t.Log(list)
}

func TestGetMail24(t *testing.T) {
	//GetMail24()
	var test string
	//ctx, cancel := Apply(true)
	//defer cancel()
	ctx, _ := Apply(true)
	err := chromedp.Run(ctx, GetMail24MailName(&test))
	t.Log(err)
	t.Log(test)
	err = chromedp.Run(ctx, GetMail24LatestMail(&test))
	t.Log(err)
	fmt.Println(test)
}

func TestApply(t *testing.T) {
	ctx, cancel := Apply(false)
	defer cancel()
	var res string
	err := chromedp.Run(ctx, AntiDetectionHeadless(), chromedp.Tasks{
		chromedp.Sleep(20 * time.Second),
		// 跳转页面
		//chromedp.Navigate("https://intoli.com/blog/not-possible-to-block-chrome-headless/chrome-headless-test.html"),
		chromedp.Navigate("https://www.pexels.com/zh-cn/new-photos?page=1"),
		// 读取HTML源码
		chromedp.InnerHTML("html", &res, chromedp.BySearch),
	})
	t.Log(err)
	t.Log(res)
	url := "https://www.pexels.com/zh-cn/photo/3584157/"
	// 新建浏览器标签页及上下文
	ctx, cancel = chromedp.NewContext(ctx, chromedp.WithTargetID(target.ID(target.CreateTarget(url).BrowserContextID)))
	defer cancel()
	err = chromedp.Run(ctx, AntiDetectionHeadless(), chromedp.Tasks{
		// 读取HTML源码
		//chromedp.InnerHTML("html", &res, chromedp.BySearch),
	})
	t.Log(err)
	t.Log(res)
}

// https://pkg.go.dev/net/mail#ReadMessage
func TestMail(t *testing.T) {
	msg := "Received: from a27-154.smtp-out.us-west-2.amazonses.com ([54.240.27.154]) by temporary-mail.net\n  for <yikxcsuka@meantinc.com>; Wed, 21 Apr 2021 13:29:21 +0800 (CST)\nDKIM-Signature: v=1; a=rsa-sha256; q=dns/txt; c=relaxed/simple;\n\ts=n6yk34xlzntpmtevqgs5ghp2jksprvft; d=netsarang.com; t=1618982959;\n\th=Date:To:From:Reply-To:Subject:Message-ID:MIME-Version:Content-Type:Content-Transfer-Encoding;\n\tbh=OsRfLaUS97/+yRoJ/BSpIARvBe+S33pKrzp1it7xCyQ=;\n\tb=LTblU9qEhTeorpht/julhD6ar7a6MDmEI9zH3TBy28KI6ah7Q+E1J0fAML2Pbcd0\n\tqX/68C8+vtGD03BGEzVPjJTJDoztO0qoNz5l7C/DSGV1MEyxh8ccQtW7rw+7+kQXv/+\n\tBKztqxpwXMjNuzZb3E0W1GYju0kr4qT5nXkMXD5o=\nDKIM-Signature: v=1; a=rsa-sha256; q=dns/txt; c=relaxed/simple;\n\ts=7v7vs6w47njt4pimodk5mmttbegzsi6n; d=amazonses.com; t=1618982959;\n\th=Date:To:From:Reply-To:Subject:Message-ID:MIME-Version:Content-Type:Content-Transfer-Encoding:Feedback-ID;\n\tbh=OsRfLaUS97/+yRoJ/BSpIARvBe+S33pKrzp1it7xCyQ=;\n\tb=NhCoNeJCoIpySzLuSxUV5P0zlsh4pLKUz5bhIG3spGhLgW0Pzf/1ZHyJJOTL9T3C\n\tDi7ChUyfnWdRjWSp+7EiA4VNrqGtOzoOsMZKipURghknjlG8bYjpCpdGXrO1D5IBlIj\n\tCl8sZevNShTh++kQW27S4S83cuDhbwxGEgxhSy3k=\nDate: Wed, 21 Apr 2021 05:29:19 +0000\nTo: yikxcsuka@meantinc.com\nFrom: \"NetSarang, Inc.\" <no-reply@netsarang.com>\nReply-To: no-reply@netsarang.com\nSubject: Xshell 7 download instruction\nMessage-ID: <01010178f2e77a87-c9bb10f6-ef08-4dad-ba57-4524546de5d2-000000@us-west-2.amazonses.com>\nX-Mailer: PHPMailer 5.2.10 (https://github.com/PHPMailer/PHPMailer/)\nMIME-Version: 1.0\nContent-Type: text/html; charset=utf-8\nContent-Transfer-Encoding: base64\nFeedback-ID: 1.us-west-2.l7ekw14vD6Jumpwas0GHbg0O54ld7FbCklw8tqJLu88=:AmazonSES\nX-SES-Outgoing: 2021.04.21-54.240.27.154\n\nPHNwYW4+RGVhciB1c2VyLDwvc3Bhbj4NCjxiciAvPjxiciAvPg0KPHNwYW4+VGhhbmsgeW91IGZv\nciB5b3VyIGludGVyZXN0IGluIFhzaGVsbCA3LiBXZSBoYXZlIHByZXBhcmVkIHlvdXIgZXZhbHVh\ndGlvbiBwYWNrYWdlLiBJZiB5b3UgZGlkIG5vdCByZXF1ZXN0IGFuIGV2YWx1YXRpb24gb2YgWHNo\nZWxsIDcsIHBsZWFzZSBjb250YWN0IG91ciBzdXBwb3J0IHRlYW0gYXQgc3VwcG9ydEBuZXRzYXJh\nbmcuY29tIHRvIGhhdmUgeW91ciBlbWFpbCBhZGRyZXNzIHJlbW92ZWQgZnJvbSBhbnkgZnV0dXJl\nIGVtYWlscyByZWxhdGVkIHRvIFhzaGVsbCA3Ljwvc3Bhbj4NCjxiciAvPjxiciAvPg0KPHNwYW4+\nUGxlYXNlIGdvIHRvIHRoZSBmb2xsb3dpbmcgVVJMIHRvIHN0YXJ0IGRvd25sb2FkaW5nIHlvdXIg\nZXZhbHVhdGlvbiBzb2Z0d2FyZTo8L3NwYW4+DQo8YnIgLz48YnIgLz4NCjxzcGFuPjxhIGhyZWY9\nImh0dHBzOi8vd3d3Lm5ldHNhcmFuZy5jb20vZW4vZG93bmxvYWRpbmcvP3Rva2VuPVZscDJUM1pN\nTVROblR6RmhXWGxFYnpWck1uSjNRVUE0U1VWWGMxcFlMVE5aY0hoRFQxcDJkV1JEU1RCQiIgdGFy\nZ2V0PSJfYmxhbmsiPmh0dHBzOi8vd3d3Lm5ldHNhcmFuZy5jb20vZW4vZG93bmxvYWRpbmcvP3Rv\na2VuPVZscDJUM1pNTVROblR6RmhXWGxFYnpWck1uSjNRVUE0U1VWWGMxcFlMVE5aY0hoRFQxcDJk\nV1JEU1RCQjwvYT48L3NwYW4+DQo8YnIgLz48YnIgLz4NCjxzcGFuPlRoaXMgbGluayB3aWxsIGV4\ncGlyZSBvbiBNYXkgMjEsIDIwMjE8L3NwYW4+IDxzcGFuPllvdSBjYW4gZXZhbHVhdGUgdGhlIHNv\nZnR3YXJlIGZvciAzMCBkYXlzIHNpbmNlIGluc3RhbGxhdGlvbi48L3NwYW4+DQo8YnIgLz48YnIg\nLz48YnIgLz4NCjxiPkRvIHlvdSBoYXZlIGFueSBxdWVzdGlvbnM/PC9iPg0KPGJyIC8+DQo8c3Bh\nbj5XZSBvZmZlciBmcmVlIHRlY2huaWNhbCBzdXBwb3J0IGR1cmluZyB0aGUgZXZhbHVhdGlvbiBw\nZXJpb2QuIElmIHlvdSBoYXZlIGFueSBxdWVzdGlvbnMsIHBsZWFzZSBzZW5kIHVzIGFuIGVtYWls\nIGF0IDxhIGhyZWY9Im1haWx0bzpzdXBwb3J0QG5ldHNhcmFuZy5jb20iPnN1cHBvcnRAbmV0c2Fy\nYW5nLmNvbTwvYT4uPC9zcGFuPg0KPGJyIC8+PGJyIC8+PGJyIC8+DQo8c3Bhbj5CZXN0IHJlZ2Fy\nZHMsPC9zcGFuPg0KPGJyIC8+PGJyIC8+PGJyIC8+DQo8dGFibGUgYm9yZGVyPSIwIiBjZWxscGFk\nZGluZz0iMCIgY2VsbHNwYWNpbmc9IjAiPg0KPHRib2R5Pg0KPHRyPjx0ZD49PT09PT09PT09PT09\nPT09PT09PT09PT09PT09PT09PT09PT09PT09PT09PT09PT09PT09PT09PT09PT09PT09PT09PC90\nZD48L3RyPg0KPHRyPjx0ZD5OZXRTYXJhbmcsIEluYy48L3RkPjwvdHI+DQo8dHI+PHRkPjQ3MDEg\nUGF0cmljayBIZW5yeSBEci4gQkxERyAyMiwgU3VpdGUgMTM3LCBTYW50YSBDbGFyYSwgQ0EgOTUw\nNTQsIFUuUy5BLjwvdGQ+PC90cj4NCjx0cj48dGQ+V2Vic2l0ZTogaHR0cDovL3d3dy5uZXRzYXJh\nbmcuY29tIHwgRW1haWw6IHN1cHBvcnRAbmV0c2FyYW5nLmNvbTwvdGQ+PC90cj4NCjx0cj48dGQ+\nUGhvbmU6ICg2NjkpIDIwNC0zMzAxPC90ZD48L3RyPg0KPC90Ym9keT4NCjwvdGFibGU+DQo=\n\n"
	r := strings.NewReader(msg)
	m, err := mail.ReadMessage(r)
	if err != nil {
		log.Fatal(err)
	}
	header := m.Header
	fmt.Println("Date:", header.Get("Date"))
	fmt.Println("From:", header.Get("From"))
	fmt.Println("To:", header.Get("To"))
	fmt.Println("Subject:", header.Get("Subject"))
	fmt.Println("Content-Transfer-Encoding:", header.Get("Content-Transfer-Encoding"))

	buf := new(bytes.Buffer) // io.ReadCloser类型转换为string
	buf.ReadFrom(m.Body)
	b := buf.String()
	fmt.Println("-------", b)
}
