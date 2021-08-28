# treeprint   
将树形结构以一种便于阅读的形式呈现, 需要树中结点实现Node接口, Node接口定义如下:    
type Node interface {     
&emsp;Id() string // 返回结点Id，要求每个结点唯一     
&emsp;Children() []Node    
&emsp;String() string // 返回结点的字符串表示，要求不存在换行符     
}    
