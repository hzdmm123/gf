// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package garray

import "sync"

type StringArray struct {
    mu           sync.RWMutex  // 互斥锁
    cap          int           // 初始化设置的数组容量
    size         int           // 初始化设置的数组大小
    array        []string      // 底层数组
}

func NewStringArray(size int, cap ... int) *StringArray {
    a     := &StringArray{}
    a.size = size
    if len(cap) > 0 {
        a.cap   = cap[0]
        a.array = make([]string, size, cap[0])
    } else {
        a.array = make([]string, size)
    }
    return a
}

// 获取指定索引的数据项, 调用方注意判断数组边界
func (a *StringArray) Get(index int) string {
    a.mu.RLock()
    value := a.array[index]
    a.mu.RUnlock()
    return value
}

// 设置指定索引的数据项, 调用方注意判断数组边界
func (a *StringArray) Set(index int, value string) {
    a.mu.Lock()
    a.array[index] = value
    a.mu.Unlock()
}

// 在当前索引位置前插入一个数据项, 调用方注意判断数组边界
func (a *StringArray) Insert(index int, value string) {
    a.mu.Lock()
    rear   := append([]string{}, a.array[index : ]...)
    a.array = append(a.array[0 : index], value)
    a.array = append(a.array, rear...)
    a.mu.Unlock()
}

// 删除指定索引的数据项, 调用方注意判断数组边界
func (a *StringArray) Remove(index int) {
    a.mu.Lock()
    a.array = append(a.array[ : index], a.array[index + 1 : ]...)
    a.mu.RUnlock()
}

// 追加数据项
func (a *StringArray) Append(value string) {
    a.mu.Lock()
    a.array = append(a.array, value)
    a.mu.Unlock()
}

// 数组长度
func (a *StringArray) Len() int {
    a.mu.RLock()
    length := len(a.array)
    a.mu.RUnlock()
    return length
}

// 返回原始数据数组
func (a *StringArray) Slice() []string {
    a.mu.RLock()
    array := a.array
    a.mu.RUnlock()
    return array
}

// 清空数据数组
func (a *StringArray) Clear() {
    a.mu.Lock()
    if a.cap > 0 {
        a.array = make([]string, a.size, a.cap)
    } else {
        a.array = make([]string, a.size)
    }
    a.mu.Unlock()
}

// 使用自定义方法执行加锁修改操作
func (a *StringArray) LockFunc(f func(array []string)) {
    a.mu.Lock()
    f(a.array)
    a.mu.Unlock()
}

// 使用自定义方法执行加锁读取操作
func (a *StringArray) RLockFunc(f func(array []string)) {
    a.mu.RLock()
    f(a.array)
    a.mu.RUnlock()
}