package sqlkit

import "github.com/mizuki1412/go-core-kit/library/commonkit"

// TxArea 事务物理代码块，不指定datasource时，用defaultDataSource
func TxArea(f func(targetDS *DataSource), dataSources ...*DataSource) {
	var ds *DataSource
	if len(dataSources) == 0 {
		ds = DefaultDataSource()
	} else {
		ds = dataSources[0]
	}
	ex := commonkit.RecoverFuncWrapper(func() {
		ds.BeginTX()
		// 传入带tx的datasource，内部代码用这个ds
		f(ds)
		ds.Commit()
	})
	if ex != nil {
		ds.Rollback()
		panic(ex)
	}
}
