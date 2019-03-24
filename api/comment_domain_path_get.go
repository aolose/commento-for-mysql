package main

import ()

func commentDomainPathGet(commentHex string) (string, string, error) {
	if commentHex == "" {
		return "", "", errorMissingField
	}
	type result struct {
		Domain string
		Path   string
	}
	var r result
	db.Table("comments").Select("domain, path").Where("comment_hex = ?", commentHex).Scan(&r)
	if r.Domain != "" {
		return r.Domain, r.Path, nil
	}
	return "", "", errorNoSuchDomain
}
