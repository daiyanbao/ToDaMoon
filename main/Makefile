# 设置编译后文件的名称
BINARY=ToDaMoon
# 编译后，把可执行文件移去FILE_PATH
FILE_PATH=`echo ${HOME}/ToDaMoon/`

# 编译号文件名称
BUILD_VERSION_FILE=Version
BUILD_NUMBER_FILE=BuildNumber
# 抓取当前git最新的hash码
GITHASH=`git rev-parse HEAD`
# 编译日期
BUILD=`date +%FT%T%z`


# 设置go程序中，对应变量的值
LDFLAGS=-ldflags "-w -s -X main.GitHash=${GITHASH} -X main.Version=$$(cat $(BUILD_VERSION_FILE)) -X main.BuildNumber=$$(cat $(BUILD_NUMBER_FILE)) -X main.BuildTime=${BUILD}"

# 编译程序
build:	
	@./versionMaker
	@go build ${LDFLAGS} -o ${BINARY}_V$$(cat $(BUILD_VERSION_FILE)).$$(cat $(BUILD_NUMBER_FILE))
	@cp ${BINARY}_V$$(cat $(BUILD_VERSION_FILE)).$$(cat $(BUILD_NUMBER_FILE)) ${BINARY}
	@rm ${BUILD_NUMBER_FILE}
	./${BINARY} version	
	mv ${BINARY}* ${FILE_PATH}
# Installs our project: copies binaries
install:
	go install ${LDFLAGS}

# Cleans our project: deletes binaries
clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi

.PHONY: clean install
