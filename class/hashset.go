package class

import (
	"errors"
	"reflect"
)

type HashSet struct {
	//数据载体
	data map[interface{}]interface{}
	//数据类型
	dataType string
	//数据数量
	count int
}

/**
初始化并指定存储对象的类型
*/
func NewHashSet(data interface{}) *HashSet {
	hashSet := new(HashSet)
	hashSet.data = make(map[interface{}]interface{})
	hashSet.dataType = reflect.TypeOf(data).String()
	return hashSet
}

/**
返回数据数量
*/
func (hashSet *HashSet) Size() int {
	return hashSet.count
}

/**
返回数据类型
*/
func (hashSet *HashSet) GetDataType() interface{} {
	return hashSet.dataType
}

/**
添加元素
*/
func (hashSet *HashSet) Add(key interface{}) error {
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

/**
删除指定Key元素
*/
func (hashSet *HashSet) Remove(key interface{}) error {
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

/**
判断key是否存在
*/
func (hashSet *HashSet) Contains(key interface{}) (bool, error) {
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

/**
重置
*/
func (hashSet *HashSet) Clear() {
	hashSet.count = 0
	hashSet.data = make(map[interface{}]interface{})
}

/**
判断添加元素是否为指定类型
*/
func (hashSet *HashSet) checkData(data interface{}) error {
	if data == nil {
		return errors.New("dataIsNil")
	}
	dataTypeof := reflect.TypeOf(data).String()
	if hashSet.dataType != dataTypeof {
		return errors.New("UnsupportedTypes")
	}
	return nil
}
