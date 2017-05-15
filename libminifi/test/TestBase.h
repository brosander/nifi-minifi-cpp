/**
 *
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

#ifndef LIBMINIFI_TEST_TESTBASE_H_
#define LIBMINIFI_TEST_TESTBASE_H_
#include <dirent.h>
#include <cstdio>
#include <cstdlib>
#include "ResourceClaim.h"
#include "catch.hpp"
#include <vector>
#include "core/logging/Logger.h"
#include "core/Core.h"
#include "properties/Configure.h"
#include "properties/Properties.h"
#include "core/logging/LoggerConfiguration.h"

class LogTestController {
 public:
  LogTestController(const std::string level = "debug") {
//     logging::Logger::getLogger()->setLogLevel(level);
   logging::LoggerProperties *logger_properties = new logging::LoggerProperties();
   logging::LoggerConfiguration::initialize(logger_properties);
   delete logger_properties;
  }

  void enableDebug() {
//     logging::Logger::getLogger()->setLogLevel("debug");
  }

  ~LogTestController() {
//     logging::Logger::getLogger()->setLogLevel(logging::LOG_LEVEL_E::info);
  }
};

class TestController {
 public:

  TestController()
      : log("info") {
    minifi::ResourceClaim::default_directory_path = const_cast<char*>("./");
  }

  ~TestController() {
    for (auto dir : directories) {
      DIR *created_dir;
      struct dirent *dir_entry;
      created_dir = opendir(dir);
      if (created_dir != NULL) {
        while ((dir_entry = readdir(created_dir)) != NULL) {
          if (dir_entry->d_name[0] != '.') {

            std::string file(dir);
            file += "/";
            file += dir_entry->d_name;
            unlink(file.c_str());
          }
        }
      }
      closedir(created_dir);
      rmdir(dir);
    }
  }

  void enableDebug(const std::string name) {
    logging::LoggerConfiguration::getConfiguration()->get_logger(name)->set_level(spdlog::level::debug);
  }

  char *createTempDirectory(char *format) {
    char *dir = mkdtemp(format);
    directories.push_back(dir);
    return dir;
  }

 protected:
  LogTestController log;
  std::vector<char*> directories;

};

#endif /* LIBMINIFI_TEST_TESTBASE_H_ */
