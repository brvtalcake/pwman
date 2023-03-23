package main

import (
	/* "golang.org/x/crypto/ssh/terminal" */
	/* crypt "pwman/src/encryption"
	encrypt "pwman/src/encryption" */

	"log"
	app "pwman/src/app"
)

func main() {
	pwman := app.PWMan_App{}
	pwman.Init()
	for !pwman.Quit {
		if !pwman.Authorized_key {
			pwman.RunEntryForm()
		} else {
			pwman.RunPswdList()
		}
	}

	/* pwman.Key = "12345678901234567890123456789012" // 32 bytes
	pwman.Byte_key = pwman.ConvertKey()
	pwman.Authorized_key = true */
	/*pwman.AddToArchive([]string{"Facebook", "554eha5ezrgx64#%$#"})
	println("Added to archive : " + pwman.Archive.DecryptedContent)
	pwman.AddToArchive([]string{"Twitter", "554eha5x64#erbvqzr$#"})
	println("Added to archive : " + pwman.Archive.DecryptedContent)
	pwman.AddToArchive([]string{"Instagram", "554eha5x6efefaze4554#%$#"})
	println("Added to archive : " + pwman.Archive.DecryptedContent) */

	/* entries := pwman.ParseArchive()
	for _, entry := range entries {
		println("Password for " + entry[0] + " is " + entry[1])
	} */

	/* os.Create("text2.bz2")
	io_writer, err := os.OpenFile("text2.bz2", os.O_WRONLY, 0755)
	if err != nil {
		panic(err.Error())
	}
	bz_writer, err := bz.NewWriter(io_writer, &bz.WriterConfig{Level: bz.BestCompression})
	if err != nil {
		panic(err.Error())
	} */
	/* defer io_writer.Close()
	defer bz_writer.Close() */
	/* io.WriteString(bz_writer, "test test test test test")
	bz_writer.Close()
	io_writer.Close()

	io_reader, err := os.OpenFile("text2.bz2", os.O_RDONLY, 0755)
	if err != nil {
		panic(err.Error())
	}
	bz_reader := cmp.NewReader(io_reader)
	if err != nil {
		panic(err.Error())
	}
	defer io_reader.Close()
	result, err := io.ReadAll(bz_reader)
	if err != nil {
		panic(err.Error())
	}
	println(string(result)) */

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

	/* d_slice := make([][]string, 0)
	fmt.Printf("%d\n", len(d_slice))
	d_slice = append(d_slice, []string{"test", "test2"})
	fmt.Printf("%d\n", len(d_slice))
	fmt.Printf("%d\n", len(d_slice[0]))
	fmt.Printf("len(%s) = %d\n", d_slice[0][0], len(d_slice[0][0]))
	fmt.Printf("len(%s) = %d\n", d_slice[0][1], len(d_slice[0][1])) */
	pwman.ClearAppResources()
	log.Println("Quiting...")
}
