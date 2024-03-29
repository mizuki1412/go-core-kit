package class

import (
	"errors"
	"reflect"
)

type HashSet struct {
	//数据载体
	data map[any]any
	//数据类型
	dataType string
	//数据数量
	count int
}

// NewHashSet 初始化并指定存储对象的类型
func NewHashSet(data any) *HashSet {
	hashSet := new(HashSet)
	hashSet.data = make(map[any]any)
	hashSet.dataType = reflect.TypeOf(data).String()
	return hashSet
}

// Size 返回数据数量
func (hashSet *HashSet) Size() int {
	return hashSet.count
}

// GetDataType 返回数据类型
func (hashSet *HashSet) GetDataType() any {
	return hashSet.dataType
}

// Add 添加元素
func (hashSet *HashSet) Add(key any) error {
	err := hashSet.checkData(key)
	if err != nil {
		return err
	}
	_, ok := hashSet.data[key]
	if ok {
		return errors.New("DataIsExist")
	}
	hashSet.count += 1
	hashSet.data[key] = key
	return nil
}

// Remove 删除指定Key元素
func (hashSet *HashSet) Remove(key any) error {
	err := hashSet.checkData(key)
	if err != nil {
		return err
	}

	value, ok := hashSet.data[key]
	if ok {
		delete(hashSet.data, value)
		hashSet.count -= 1
		return nil
	}
	return errors.New("NotFoundKey")
}

// Contains 判断key是否存在
func (hashSet *HashSet) Contains(key any) (bool, error) {
	err := hashSet.checkData(key)
	if err != nil {
		return false, err
	}
	_, ok := hashSet.data[key]
	if ok {
		return true, nil
	} else {
		return false, nil
	}
}

// Clear 重置
func (hashSet *HashSet) Clear() {
	hashSet.count = 0
	hashSet.data = make(map[any]any)
}

//判断添加元素是否为指定类型
func (hashSet *HashSet) checkData(data any) error {
	if data == nil {
		return errors.New("dataIsNil")
	}
	dataTypeof := reflect.TypeOf(data).String()
	if hashSet.dataType != dataTypeof {
		return errors.New("UnsupportedTypes")
	}
	return nil
}
