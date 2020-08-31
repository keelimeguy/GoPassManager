我是光年实验室高级招聘经理。
我在github上访问了你的开源项目，你的代码超赞。你最近有没有在看工作机会，我们在招软件开发工程师，拉钩和BOSS等招聘网站也发布了相关岗位，有公司和职位的详细信息。
我们公司在杭州，业务主要做流量增长，是很多大型互联网公司的流量顾问。公司弹性工作制，福利齐全，发展潜力大，良好的办公环境和学习氛围。
公司官网是http://www.gnlab.com,公司地址是杭州市西湖区古墩路紫金广场B座，若你感兴趣，欢迎与我联系，
电话是0571-88839161，手机号：18668131388，微信号：echo 'bGhsaGxoMTEyNAo='|base64 -D ,静待佳音。如有打扰，还请见谅，祝生活愉快工作顺利。

# GoPassManager
A simple password manager written in go.

## Dependencies:
PostgreSQL + pgcrypto

github.com/lib/pq/commit/90697d60dd844d5ef6ff15135d0203f65d2f53b8

## Built On:
Windows Subsystem for Linux (Debian) on Windows 10 Version 1803

`cat /etc/issue` : Debian GNU/Linux 9 \n \l

`uname -a` : Linux ~~REDACTED~~ 4.4.0-17134-Microsoft #137-Microsoft Thu Jun 14 18:46:00 PST 2018 x86_64 GNU/Linux

`go version` : go version go1.10.2 linux/amd64

`bash --version` : GNU bash, version 4.4.12(1)-release (x86_64-pc-linux-gnu)

`psql --version` : psql (PostgreSQL) 9.6.7

`gpg --version` : gpg (GnuPG) 2.1.18 \n libgcrypt 1.7.6-beta

`echo $TERM` : xterm-256color

## Security points:
This is mainly a personal project for the fun of it. I am not a security expert.

Take your own precautions and use at your own discretion.

Potential issues:
1. Any security issues that may exist in PGP (RFC4880) and your related tools
2. Any security isses that may exist in PostgreSQL, including pgcrypto and pgp_pub_encrypt/decrypt(..)
4. Potential logging of stdout in your terminal, including within `tput smcup` environment
5. Potential internal logging within `/bin/bash`
6. Passwords sitting in memory that is accessible to other processes
7. User mistakes, including leaving the program open or allowing the executable to be replaced with a malicious copy

## Building and running code:
`make && ./run`
