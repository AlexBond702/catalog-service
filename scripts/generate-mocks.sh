#!/bin/bash

# Параметры mockery:
#   --dir - директория, где находится интерфейс для генерации мока
#   --name - имя интерфейса, для которого нужно создать мок
#   --output - директория, куда будет записан сгенерированный мок
#   --filename - имя файла для сгенерированного мока
#   --outpkg - имя пакета в сгенерированном моке
#   --structname - имя структуры мока (если не указано, используется имя интерфейса + "Mock")
#   --disable-version-string - не добавлять версию mockery в комментарии сгенерированного файла
#   --with-expecter - использовать современный EXPECT() API вместо устаревшего On()

# Мок для product.Repository
mockery \
  --dir=internal/app/repository/ \
  --name=Product \
  --output=internal/app/repository/product \
  --filename=repository_mock.go \
  --outpkg=pproduct \
  --structname=MockProduct \
  --disable-version-string \
  --with-expecter

# Мок для category.Repository
mockery \
  --dir=internal/app/repository/ \
  --name=Category \
  --output=internal/app/repository/category \
  --filename=repository_mock.go \
  --outpkg=pproduct \
  --structname=MockCategory \
  --disable-version-string \
  --with-expecter


