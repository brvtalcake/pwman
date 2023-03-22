package main

import (
	/* "golang.org/x/crypto/ssh/terminal" */
	/* crypt "pwman/src/encryption"
	encrypt "pwman/src/encryption" */

	app "pwman/src/app"
)

func main() {
	pwman := app.PWMan_App{}
	pwman.Init()
	for !pwman.Quit {
		for !pwman.Authorized_key {
			pwman.RunEntryForm()
		}
	}
	/* println("BZ2 test") */
	/* os.Create("test.bz2")
	io_writer, err := os.OpenFile("test.bz2", os.O_RDWR, 0755)
	if err != nil {
		println(err.Error())
	}
	defer io_writer.Close()
	writter, err := bz.NewWriter(io_writer, &bz.WriterConfig{Level: bz.BestCompression})
	if err != nil {
		println(err.Error())
	}
	defer writter.Close()
	writter.Write([]byte("Hello world ! Lorem ipsum dolor sit amet, consectetur adipiscing elit. Donec auctor, nisl nec ultricies luctus, nisl nisl luctus nisl, nec luctus nisl nisl nec nisl. Donec auctor, nisl nec ultricies luctus, nisl nisl luctus nisl, nec luctus nisl nisl nec nisl. Donec auctor, nisl nec ultricies luctus, nisl nisl luctus nisl, nec luctus nisl nisl nec nisl. Donec auctor, nisl nec ultricies luctus, nisl nisl luctus nisl, nec luctus nisl nisl nec nisl. Donec auctor, nisl nec ultricies luctus, nisl nisl luctus nisl, nec luctus nisl nisl nec nisl. Donec auctor, nisl nec ultricies luctus, nisl nisl luctus nisl, nec luctus nisl nisl nec nisl. Donec auctor, nisl nec ultricies luctus, nisl nisl luctus nisl, nec luctus nisl nisl nec nisl. Donec auctor, nisl nec ultricies luctus, nisl nisl luctus nisl, nec luctus nisl nisl nec nisl. Donec auctor, nisl nec ultricies luctus, nisl nisl luctus nisl, nec luctus nisl nisl nec nisl. Donec auctor, nisl nec ultricies luctus, nisl nisl luctus nisl, nec luctus nisl nisl nec nisl. Donec auctor, nisl nec ultricies luctus, nisl nisl luctus nisl, nec luctus nisl nisl nec nisl. Donec auctor, nisl nec ultricies luctus, nisl nisl luctus nisl, nec luctus nisl nisl nec nisl. Donec auctor, nisl nec ultricies luctus, nisl nisl luctus nisl, nec luctus nisl nisl nec nisl. Donec auctor, nisl nec ultricies luctus, nisl nisl luctus nisl")) */

	/* io_reader, err := os.OpenFile("test.bz2", os.O_RDONLY, 0755)
	if err != nil {
		println(err.Error())
	}
	defer io_reader.Close()
	read_conf := bz.ReaderConfig{}
	reader, err := bz.NewReader(io_reader, &read_conf)
	if err != nil {
		println(err.Error())
	}
	defer reader.Close()
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		println(err.Error())
	}
	println(string(data))
	println("BZ2 test end")
	println("Quiting...") */
	/* pwman.RunEntryForm() */

	pwman.ClearAppResources()
	println("Quiting...")
}
