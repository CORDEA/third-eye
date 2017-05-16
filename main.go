/*
 * Copyright 2017 Yoshihiro Tanaka
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * Author: Yoshihiro Tanaka <contact@cordea.jp>
 * date  : 2017-05-16
 */

package main

import (
	"flag"
	"fmt"
	"github.com/golang/glog"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var (
	isDecrypt = flag.Bool("d", false, "")
	key       = flag.String("key", "", "")
	thirdEye  ThirdEye
)

func walk(path string, info os.FileInfo, err error) error {
	if err != nil {
		glog.Fatalln(err)
	}
	if info.IsDir() {
		return nil
	}

	if err := thirdEye.Rename(path); err != nil {
		glog.Warningln(err)
	}

	return nil
}

func renameFiles(path string) {
	stat, err := os.Stat(path)
	if err != nil {
		glog.Fatalln(err)
	}

	switch mode := stat.Mode(); {
	case mode.IsRegular():
		if err := thirdEye.Rename(path); err != nil {
			glog.Fatalln(err)
		}
		return
	case mode.IsDir():
		filepath.Walk(path, walk)
		return
	}
}

func main() {
	flag.Parse()

	if flag.NArg() == 0 {
		glog.Fatalln("Required parameters are missing.")
	}

	if aLen := len(*key); aLen != 32 {
		bytes, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			glog.Fatalln(err)
		}
		k := strings.Trim(string(bytes), "\n")
		if sLen := len(k); sLen != 32 {
			kLen := aLen
			if sLen > aLen {
				kLen = sLen
			}
			glog.Fatalln(fmt.Sprintf("Key length is wrong. actual: %d, expected: 32.", kLen))
		}
		key = &k
	}

	thirdEye = ThirdEye{
		Key:       *key,
		IsDecrypt: *isDecrypt,
	}
	renameFiles(flag.Arg(0))
}
