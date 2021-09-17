#!/usr/bin/python3

import os
import json
from mutagen.mp3 import MP3

# Run this script befor ampgosetup.  This provides the 
# duration of the mp3 in milliseconds, to be used in the 
# players slider.

class FindDuration:

    def __init__(self):
        self.mp3_path = "/media/pipi/FOO/music"

    def convert_to_json(self, alist):
        jsonStr = json.dumps(alist)
        return jsonStr

    def write_to_file(self, jsonstr, ofile):
        with open(ofile, r"w") as newfile:
            newfile.write(jsonstr)

    def mp3_info(self, afile):
        audio = MP3(afile)
        audio_info = audio.info
        length_in_secs = int(audio_info.length)
        length_in_mills = length_in_secs * 1000
        return length_in_mills

    def find_files(self, apath):
        for (paths, dirs, files) in os.walk(apath):
            for filename in files:
                print("Processing:\n %s" % filename)
                fnn = os.path.join(paths, filename)
                fpath, ext = os.path.splitext(fnn)
                outfile = fpath + ".mp3info"
                newfn = fpath.replace("/media/pipi/FOO/music", "/root/fsData")
                newfilename = newfn + ext
                if ext == ".mp3":
                    duration = self.mp3_info(fnn)
                    x = {}
                    x['filename'] = newfilename
                    x['duration'] = str(duration)
                    jstring = self.convert_to_json(x)
                    self.write_to_file(jstring, outfile)

    def remove_old(self, apath):
        for (paths, dirs, files) in os.walk(apath):
            for filename in files:
                print("Processing:\n %s" % filename)
                fnn = os.path.join(paths, filename)
                ext = os.path.splitext(fnn)[1]
                if ext == ".mp3info":
                    os.remove(fnn)

    def main(self):
        self.remove_old(self.mp3_path)
        mp3list = self.find_files(self.mp3_path)

if __name__ == "__main__":
    fd = FindDuration()
    fd.main()