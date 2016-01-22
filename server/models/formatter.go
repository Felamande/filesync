package models

type Formatter interface{
    FormatTo(to interface{})error
}