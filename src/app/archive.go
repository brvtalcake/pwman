package app

import (
	"io"
	"log"
	"os"
	"path"
	"pwman/src/encryption"
	"strings"

	bz "github.com/dsnet/compress/bzip2"
)

const archive_name string = "pwman_archive.bz2"

func LookForArchive() (bool, error) {
	exec_path, err := os.Executable()
	if err != nil {
		return false, err
	}
	exec_path = path.Dir(exec_path)
	archive_path := exec_path + "/" + archive_name
	return CheckIfExists(archive_path)
}

func CheckIfExists(file_name string) (bool, error) {
	if _, err := os.Stat(file_name); !os.IsNotExist(err) {
		return true, err
	}
	return false, nil
}

func (ar *PWMan_Archive) GetPhysicalArchivePath() (string, error) {
	if exec_path, err := os.Executable(); err != nil {
		ar.Location = ""
		return "", err
	} else {
		exec_path = path.Dir(exec_path)
		ar.Location = exec_path + "/" + archive_name
		return exec_path + "/" + archive_name, nil
	}
}

func (this_app *PWMan_App) CreateHeader() *PWMan_App {
	this_app.AddToArchive([]string{"PWMAN_ARCHIVE", "0.0.1"})
	return this_app
}

func (this_app *PWMan_App) ParseArchive() [][]string {
	var err error
	var returned_entries [][]string
	// slurp the whole encrypted content and store it in archive struct
	if this_app.Archive.BZ2_Writer != nil {
		err = this_app.Archive.BZ2_Writer.Close()
		if err != nil {
			log.Panic(err.Error())
			this_app.ClearAppResources()
			os.Exit(1)
		}
	}
	if this_app.Archive.BZ2_Reader == nil {
		this_app.Archive.BZ2_Reader, err = bz.NewReader(this_app.Archive.IO_Reader, &bz.ReaderConfig{})
		if err != nil {
			log.Panic(err.Error())
			this_app.ClearAppResources()
			os.Exit(1)
		}
	} else if this_app.Archive.BZ2_Reader != nil { // close and reopen the reader to be sure
		err = this_app.Archive.BZ2_Reader.Close()
		if err != nil {
			log.Panic(err.Error())
			this_app.ClearAppResources()
			os.Exit(1)
		}
		this_app.Archive.BZ2_Reader, err = bz.NewReader(this_app.Archive.IO_Reader, &bz.ReaderConfig{})
		if err != nil {
			log.Panic(err.Error())
			this_app.ClearAppResources()
			os.Exit(1)
		}
	}
	this_app.Archive.BZ2_Reader.Reset(this_app.Archive.IO_Reader)
	_, err = this_app.Archive.IO_Reader.Seek(0, io.SeekStart)
	if err != nil {
		log.Panic(err.Error())
		this_app.ClearAppResources()
		os.Exit(1)
	}
	this_app.Archive.EncryptedContent, err = io.ReadAll(this_app.Archive.BZ2_Reader) // maybe to change
	if err != nil {
		log.Panic(err.Error())
		this_app.ClearAppResources()
		os.Exit(1)
	}
	this_app.Archive.DecryptedContent, err = encryption.Decrypt(this_app.Archive.EncryptedContent, this_app.Byte_key)
	if err != nil {
		if err.Error() == "cipher: message authentication failed" {
			log.Panic("Wrong key !")
			log.Println("\n\x1b[34;1mError :\n\x1b[0m" + err.Error())
			recover()
		} else {
			log.Panic(err.Error())
			this_app.ClearAppResources()
			os.Exit(1)
		}
	}

	runed_decrypted_content := []rune(this_app.Archive.DecryptedContent)
	runed_decrypted_content_len := len(runed_decrypted_content)
	no_more_entries := false
	global_cursor := uint64(0)
	entry_count := uint64(0)

	check_if_no_more_entries := func() bool {
		return global_cursor >= uint64(runed_decrypted_content_len)-1
	}
	get_next_entry := func() []string {
		next_entry := make([]string, 2)
		no_more_entries = check_if_no_more_entries()
		for !IsValidEntryChar(runed_decrypted_content[global_cursor]) && !no_more_entries {
			global_cursor++
			no_more_entries = check_if_no_more_entries()
		}
		for IsValidEntryChar(runed_decrypted_content[global_cursor]) && !no_more_entries {
			next_entry[0] += string(runed_decrypted_content[global_cursor])
			global_cursor++
			no_more_entries = check_if_no_more_entries()
		}
		for !IsValidEntryChar(runed_decrypted_content[global_cursor]) && !no_more_entries {
			global_cursor++
			no_more_entries = check_if_no_more_entries()
		}
		for IsValidEntryChar(runed_decrypted_content[global_cursor]) && !no_more_entries {
			next_entry[1] += string(runed_decrypted_content[global_cursor])
			global_cursor++
			no_more_entries = check_if_no_more_entries()
		}
		entry_count++
		if !strings.Contains(next_entry[0], "PWMAN_ARCHIVE") && !strings.Contains(next_entry[1], "PWMAN_ARCHIVE") { // skip the header
			return next_entry
		} else {
			return nil
		}

	}

	// split the decrypted content into entries
	/*
	* Each entry has the following format :
	* 1.	Name of the entry (website, app, etc.)
	* 2.	Associated password
	* 3.	2 newlines
	 */

	for global_cursor < uint64(runed_decrypted_content_len) && !no_more_entries {
		returned_entries = append(returned_entries, get_next_entry())
	}
	this_app.Archive.PswEntryCount = entry_count - 1
	err = this_app.Archive.BZ2_Reader.Close()
	if err != nil {
		log.Panic(err.Error())
		this_app.ClearAppResources()
		os.Exit(1)
	}
	if entry_count > 1 { // skip the header
		return ClearVoidEntries(returned_entries)
	} else {
		return nil
	}
}

func ClearVoidEntries(entries [][]string) [][]string {
	var returned_entries [][]string = nil
	for _, entry := range entries {
		if entry != nil && entry[0] != "" && entry[1] != "" {
			returned_entries = append(returned_entries, entry)
		}
	}
	return returned_entries
}

/* func (this_app *PWMan_App) ReshapeArchive() { */

func (this_app *PWMan_App) AddToArchive(entry []string) {
	var err error
	if this_app.Archive.BZ2_Reader != nil {
		err = this_app.Archive.BZ2_Reader.Close()
		if err != nil {
			log.Panic(err.Error())
			this_app.ClearAppResources()
			os.Exit(1)
		}
	}
	if this_app.Archive.BZ2_Writer == nil {
		err = os.Truncate(this_app.Archive.Location, 0)
		if err != nil {
			log.Panic(err.Error())
			this_app.ClearAppResources()
			os.Exit(1)
		}
		this_app.Archive.BZ2_Writer, err = bz.NewWriter(this_app.Archive.IO_Writer, &bz.WriterConfig{Level: bz.BestCompression})
		if err != nil {
			log.Panic(err.Error())
			this_app.ClearAppResources()
			os.Exit(1)
		}
	} else if this_app.Archive.BZ2_Writer != nil {
		err = this_app.Archive.BZ2_Writer.Close()
		if err != nil {
			log.Panic(err.Error())
			this_app.ClearAppResources()
			os.Exit(1)
		}
		err = os.Truncate(this_app.Archive.Location, 0)
		if err != nil {
			log.Panic(err.Error())
			this_app.ClearAppResources()
			os.Exit(1)
		}
		this_app.Archive.BZ2_Writer, err = bz.NewWriter(this_app.Archive.IO_Writer, &bz.WriterConfig{Level: bz.BestCompression})
		if err != nil {
			log.Panic(err.Error())
			this_app.ClearAppResources()
			os.Exit(1)
		}
	}

	this_app.Archive.BZ2_Writer.Reset(this_app.Archive.IO_Writer)
	_, err = this_app.Archive.IO_Writer.Seek(0, io.SeekStart)
	if err != nil {
		log.Panic(err.Error())
		this_app.ClearAppResources()
		os.Exit(1)
	}

	var new_entry string

	if entry == nil {
		goto just_rewrite
	}

	new_entry = "\n" + entry[0] + "\n" + entry[1] + "\n"
	this_app.Archive.DecryptedContent += new_entry
	this_app.Archive.EncryptedContent, err = encryption.Encrypt([]byte(this_app.Archive.DecryptedContent), this_app.Byte_key)
	if err != nil {
		log.Panic(err.Error())
		this_app.ClearAppResources()
		os.Exit(1)
	}

	_, err = io.WriteString(this_app.Archive.BZ2_Writer, string(this_app.Archive.EncryptedContent))
	if err != nil {
		log.Panic(err.Error())
		this_app.ClearAppResources()
		os.Exit(1)
	}
	this_app.Archive.PswEntryCount++
	err = this_app.Archive.BZ2_Writer.Close()
	if err != nil {
		log.Panic(err.Error())
		this_app.ClearAppResources()
		os.Exit(1)
	}
	return

just_rewrite:
	_, err = io.WriteString(this_app.Archive.BZ2_Writer, string(this_app.Archive.EncryptedContent))
	if err != nil {
		log.Panic(err.Error())
		this_app.ClearAppResources()
		os.Exit(1)
	}
	err = this_app.Archive.BZ2_Writer.Close()
	if err != nil {
		log.Panic(err.Error())
		this_app.ClearAppResources()
		os.Exit(1)
	}
}

func DeleteAllSubStr(s []rune, substr []rune) (string, uint64) {
	var result string
	var found_substrings uint64 = 0
	var i int
	for i < len(s) {
		if i+len(substr) <= len(s) && string(s[i:i+len(substr)]) == string(substr) {
			i += len(substr)
			found_substrings++
		} else {
			result += string(s[i])
			i++
		}
	}
	return result, found_substrings
}

func (this_app *PWMan_App) DeleteFromArchive(entry string) bool {
	var found bool = false
	var found_substrings uint64 = 0
	if strings.Contains(this_app.Archive.DecryptedContent, entry) {
		found = true
		this_app.Archive.DecryptedContent, found_substrings = DeleteAllSubStr([]rune(this_app.Archive.DecryptedContent), []rune(entry))
		this_app.Archive.EncryptedContent, _ = encryption.Encrypt([]byte(this_app.Archive.DecryptedContent), this_app.Byte_key)
		this_app.Archive.PswEntryCount -= found_substrings
		this_app.AddToArchive(nil)
	} else {
		return found
	}
	return found
}
