package memcached

type NotStoredError string

func (nse NotStoredError) Error() string {
	return "存储错误：" + string(nse)
}

type ExistsError string

func (ee ExistsError) Error() string {
	return "存储错误：" + string(ee)
}

type NotFoundError string

func (nfe NotFoundError) Error() string {
	return "存储错误：" + string(nfe)
}

type NotValidRespError string

func (nvre NotValidRespError) Error() string {
	return "数据获取错误：" + string(nvre)
}
