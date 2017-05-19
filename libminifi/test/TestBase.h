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
#include <sstream>
#include "ResourceClaim.h"
#include "catch.hpp"
#include <vector>
#include "core/logging/Logger.h"
#include "core/Core.h"
#include "properties/Configure.h"
#include "properties/Properties.h"
#include "core/logging/LoggerConfiguration.h"
#include "spdlog/sinks/ostream_sink.h"
#include "spdlog/sinks/dist_sink.h"

class LogTestController {
 public:
  static LogTestController& getInstance() {
   static LogTestController instance;
   return instance;
  }
  
  template<typename T>
  void setDebug() {
    setLevel<T>(spdlog::level::debug);
  }
  
  template<typename T>
  void setInfo() {
    setLevel<T>(spdlog::level::info);
  }
  
  template<typename T>
  void setLevel(spdlog::level::level_enum level) {
    std::string name = core::getClassName<T>();
    modified_loggers.push_back(name);
    setLevel(name, level);
  }
  
  bool contains(std::string ending) {
    logger_->info("Looking for %s in %s.", ending, log_output.str());
    return (ending.length() > 0 && log_output.str().find(ending) != std::string::npos);
  }
  
  void reset() {
    for (auto const & name : modified_loggers) {
      setLevel(name, spdlog::level::err);
    }
    modified_loggers = std::vector<std::string>();
    log_output.str("");
    log_output.clear();
  }
  
  std::ostringstream log_output;
  
 private:
   class TestBootstrapLogger: public logging::Logger {
    public:
      TestBootstrapLogger(std::shared_ptr<spdlog::logger> logger):Logger(logger){};
   };
   class TestLoggerProperties : public logging::LoggerProperties {
    public:
     TestLoggerProperties(logging::Logger & logger) : LoggerProperties(logger) {
      set("logger.root", "ERROR,ostream");
     }
   };
  LogTestController() {
   logger_ = std::make_shared<spdlog::logger>("test main", spdlog::sinks::stderr_sink_mt::instance());
   logger_->set_level(spdlog::level::info);
   TestBootstrapLogger testBootstrapLogger(logger_);
   std::shared_ptr<logging::LoggerProperties> logger_properties = std::shared_ptr<logging::LoggerProperties>(new TestLoggerProperties(testBootstrapLogger));
   std::shared_ptr<spdlog::sinks::dist_sink_mt> dist_sink = std::make_shared<spdlog::sinks::dist_sink_mt>();
   dist_sink->add_sink(std::make_shared<spdlog::sinks::ostream_sink_mt>(log_output, true));
   dist_sink->add_sink(spdlog::sinks::stderr_sink_mt::instance());
   logger_properties->add_sink("ostream", dist_sink);
   logging::LoggerConfiguration::initialize(logger_properties);
  }
  LogTestController(LogTestController const&);
  LogTestController& operator=(LogTestController const&);
  ~LogTestController() {};

  void setLevel(const std::string name, spdlog::level::level_enum level) {
    logger_->info("Setting log level for %s to %s", name, spdlog::level::to_str(level));
    logging::LoggerConfiguration::getConfiguration()->get_logger(name)->set_level(level);
  }
  std::vector<std::string> modified_loggers;
  std::shared_ptr<spdlog::logger> logger_;
};

class TestController {
 public:

  TestController()
      : log(LogTestController::getInstance()) {
    minifi::ResourceClaim::default_directory_path = const_cast<char*>("./");
    log.reset();
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

  char *createTempDirectory(char *format) {
    char *dir = mkdtemp(format);
    directories.push_back(dir);
    return dir;
  }

 protected:
  LogTestController & log;
  std::vector<char*> directories;

};

#endif /* LIBMINIFI_TEST_TESTBASE_H_ */
