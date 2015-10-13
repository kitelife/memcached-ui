package memcached

type NotStoredError string

func (nse NotStoredError) Error() string {
	return "存储错误：" + string(nse)
}

type ExistsError string

func (ee ExistsError) Error() string {
	return "存储错误：" + string(ee)
}

type NotValidRespError string

func (nvre NotValidRespError) Error() string {
	return "发生错误：" + string(nvre)
}

type NotFoundError string

func (nfe NotFoundError) Error() string {
	return "未找到对应的键值：" + string(nfe)
}
