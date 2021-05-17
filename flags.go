package main

type stringArrayFlags []string

func (i *stringArrayFlags) String() string {
	return "String array flag"
}

func (i *stringArrayFlags) Set(s string) error {
	*i = append(*i, s)
	return nil
}
