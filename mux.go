/*
like mux router
use regexp
sam
*/

package gow

import (
	"net/url"
	"regexp"
	"strings"
)

type routerPathInfo struct {
	Path     string
	fullPath string
	handlers HandlersChain
	params   *Params
	tsr      bool
	nType    nodeType
}

// RouterPath defines a routerPathInfo array.
type RouterPath []routerPathInfo

// getMuxValue Returns the handle registered with the given path (key). The values of
// wildcards are saved to a map.
// If no handle can be found, a TSR (trailing slash redirect) recommendation is
// made if a handle exists with an extra (without the) trailing slash for the
// given path.
func (n *node) getMuxValue(path string, params *Params, unescape bool) (value nodeValue) {
	rpm := getNodeRouterPathMap(n)
	routerPath, ok := getMatchPath(path, rpm, unescape)
	if ok {
		if params != nil {
			value.params = params
			value.params = routerPath.params
		}
		value.handlers = routerPath.handlers
		value.fullPath = routerPath.fullPath
		return
	}
	return
}

// getMatchPath return routerPathInfo
// regexp match
func getMatchPath(path string, rp RouterPath, unescape bool) (*routerPathInfo, bool) {
	lastChar := path[len(path)-1:]
	if path != "/" && lastChar == "/" && !strings.Contains(path, ".") {
		path = path[:len(path)-1]
	}
	path = strings.ReplaceAll(path, "//", "/")
	//request path ignore  case
	path = strings.ToLower(path)
	debugPrint("path:%s", path)
	for _, p := range rp {
		regPath, keys := mathPath(p.Path)
		if path == regPath {
			return &p, true
		} else {
			// all reg match
			ok, _ := regexp.MatchString("^"+regPath+"$", path)
			if ok {
				valueRegexp := regexp.MustCompile(regPath)
				if unescape {
					if v, err := url.QueryUnescape(path); err == nil {
						path = v
					}
				}
				values := valueRegexp.FindStringSubmatch(path)
				params := new(Params)
				for i, k := range keys {
					*params = append(*params, Param{
						Key:   strings.ReplaceAll(strings.ReplaceAll(k, "{", ""), "}", ""),
						Value: values[i+1],
					})
				}
				p.params = params
				return &p, ok
			}
		}
	}
	return nil, false
}

var (
	intRegexp  = []byte(`(\d+)`)
	charRegexp = []byte(`(\w+)`)
	starRegexp = []byte(`(.*)`)
)

func mathPath(path string) (regPath string, keys []string) {
	var (
		replaceRegexp []byte
		nPath         string
	)
	nSplit := strings.Split(path, "/")
	// like {uid} or {uid:int}
	// replace {uid} to regexp
	replaceRegexp = charRegexp
	wildcardRegexp := regexp.MustCompile(`{\w+}`)
	for _, n := range nSplit {
		if strings.Contains(n, "{") || strings.Contains(n, "*") {
			// math {uid:int}
			if strings.Contains(n, ":int") {
				n = strings.ReplaceAll(n, ":int", "")
				replaceRegexp = intRegexp
			}

			// static /static/*filepath
			if strings.Contains(n, "*filepath") {
				n = strings.ReplaceAll(n, "*filepath", string(starRegexp))
				debugPrint("filepath:%s", n)
				replaceRegexp = starRegexp
			}

			// /*action
			if strings.Contains(n, "*action") {
				n = strings.ReplaceAll(n, "*action", string(starRegexp))
				replaceRegexp = starRegexp
			}
			key := wildcardRegexp.FindAllString(n, -1)
			keys = append(keys, key...)
			nPath = string(wildcardRegexp.ReplaceAll([]byte(n), replaceRegexp))

		} else {
			nPath = n
		}
		regPath = regPath + nPath + "/"
	}

	//去掉group时可能产生的//
	if strings.Contains(regPath, "//") {
		regPath = strings.ReplaceAll(regPath, "//", "/")
	}
	if regPath != "/" {
		regPath = regPath[:len(regPath)-1]
	}
	return regPath, keys
}

func getNodeRouterPathMap(n *node) (rp RouterPath) {
	rp = getNodeRouterPath("", rp, n)
	return
}

func getNodeRouterPath(path string, rp RouterPath, root *node) RouterPath {
	path += root.path
	if len(root.handlers) > 0 {
		rp = append(rp, routerPathInfo{
			Path:     strings.ToLower(path),
			fullPath: root.fullPath,
			handlers: root.handlers,
			nType:    root.nType,
		})
	}
	for _, child := range root.children {
		rp = getNodeRouterPath(path, rp, child)
	}
	return rp
}
