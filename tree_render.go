package tree_render

import "bytes"

type Node interface {
	Id() string // 返回结点Id，要求每个结点唯一
	Children() []Node
	String() string // 返回结点的字符串表示，要求不存在换行符
}

// 生成以root为根的前缀树的可视化表示
func Render(root Node, minLeafDistance uint32) string {
	t := &tree{
		pos:             make(map[string]position),
		children:        make(map[string][]Node),
		parent:          make(map[string]Node),
		minLeafDistance: minLeafDistance,
	}
	t.dfs(root, 0)

	buf := &bytes.Buffer{}
	curLevelNodes := []Node{root} // 当前层的结点
	for len(curLevelNodes) > 0 {
		// 写入当前层结点
		firstEmptyColumn := 0 // 当前行第1个为空的结点所在的列
		for _, v := range curLevelNodes {
			vPos := t.pos[v.Id()]
			writeN(buf, " ", vPos.rootColumn-firstEmptyColumn)
			buf.WriteString(v.String())
			firstEmptyColumn = vPos.rootColumn + len(v.String())
		}

		// 获取下一层的结点
		var nextLevelNodes []Node
		for _, v := range curLevelNodes {
			for _, u := range t.children[v.Id()] {
				nextLevelNodes = append(nextLevelNodes, u)
			}
		}

		if len(nextLevelNodes) == 0 {
			return buf.String()
		}

		buf.WriteString("\n")

		// 写入当前层结点每个结点下边的"|"，对于没有孩子结点的当前层结点不需要写入
		firstEmptyColumn = 0
		for _, v := range curLevelNodes {
			vPos := t.pos[v.Id()]
			writeN(buf, " ", vPos.rootColumn-firstEmptyColumn)
			firstEmptyColumn = vPos.rootColumn
			if len(t.children[v.Id()]) > 0 {
				buf.WriteString("|")
				firstEmptyColumn++
			}
		}
		buf.WriteString("\n")

		curLevelNodes = nextLevelNodes

		// 同一父结点的孩子写入连续的横线
		firstEmptyColumn = 0
		var lastNodeParent Node
		for _, v := range curLevelNodes {
			vPos := t.pos[v.Id()]
			vParent := t.parent[v.Id()]

			var s string
			if lastNodeParent != nil && vParent != nil && lastNodeParent.Id() == vParent.Id() {
				s = "-"
			} else {
				s = " "
			}

			writeN(buf, s, vPos.rootColumn-firstEmptyColumn)
			if len(t.children[vParent.Id()]) == 1 {
				buf.WriteString("|")
			} else {
				buf.WriteString("-")
			}
			firstEmptyColumn = vPos.rootColumn + 1
			lastNodeParent = vParent
		}

		buf.WriteString("\n")

		// 写入孩子结点上方的"|"
		firstEmptyColumn = 0
		for _, v := range curLevelNodes {
			vPos := t.pos[v.Id()]
			writeN(buf, " ", vPos.rootColumn-firstEmptyColumn)
			buf.WriteString("|")
			firstEmptyColumn = vPos.rootColumn + 1
		}

		buf.WriteString("\n")
	}

	return buf.String()
}

type position struct {
	rootColumn int // 当前结点所在列序号（列序号从0开始）
	maxColumn  int // 以当前结点为根的子树中列序号最大的结点的列序号
	minColumn  int // 以当前结点为根的子树中列序号最小的结点的列序号
}

type tree struct {
	pos             map[string]position
	children        map[string][]Node
	parent          map[string]Node
	minLeafDistance uint32
}

func (t *tree) dfs(root Node, minColumn int) {
	children, inMap := t.children[root.Id()]
	if inMap {
		panic("found the backward side, this is not a tree")
	}

	children = root.Children()
	t.children[root.Id()] = children
	if len(children) == 0 {
		t.pos[root.Id()] = position{
			rootColumn: minColumn,
			maxColumn:  minColumn + len(root.String()) - 1,
			minColumn:  minColumn,
		}
		return
	}

	oldMinColumn := minColumn
	for _, v := range children {
		t.parent[v.Id()] = root
		t.dfs(v, minColumn)
		minColumn = t.pos[v.Id()].maxColumn + 1 + int(t.minLeafDistance)
	}

	// 计算root的位置
	rootColumn := (t.pos[children[0].Id()].rootColumn + t.pos[children[len(children)-1].Id()].rootColumn) / 2
	t.pos[root.Id()] = position{
		rootColumn: rootColumn,
		maxColumn: func() int {
			childrenMaxColumn := t.pos[children[len(children)-1].Id()].maxColumn
			rootMaxColumn := rootColumn + len(root.String()) - 1
			if childrenMaxColumn < rootMaxColumn {
				return rootMaxColumn
			}
			return childrenMaxColumn
		}(),
		minColumn: oldMinColumn,
	}
}

// 向buf中写入n个str
func writeN(buf *bytes.Buffer, str string, n int) {
	for i := 1; i <= n; i++ {
		buf.WriteString(str)
	}
}
