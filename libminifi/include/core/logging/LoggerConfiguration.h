/**
 * @file LoggerConfiguration.h
 * Logger class declaration
 * This is a C++ wrapper for spdlog, a lightweight C++ logging library
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
#ifndef __LOGGER_CONFIGURATION_H__
#define __LOGGER_CONFIGURATION_H__

#include <map>
#include <string>
#include "core/logging/Logger.h"
#include "properties/Properties.h"

#include "spdlog/spdlog.h"

namespace org {
namespace apache {
namespace nifi {
namespace minifi {
namespace core {
namespace logging {
 
static const char *appender_prefix;
static const char *logger_prefix;

struct LoggerNamespace {
 spdlog::level::level_enum level;
 std::vector<std::shared_ptr<spdlog::sinks::sink>> sinks;
 std::map<std::string, std::shared_ptr<LoggerNamespace>> children;
};
 
class LoggerProperties : public Properties {
public:
 LoggerProperties(Logger & logger) : Properties(logger) {};
 std::vector<std::string> get_appenders();
 std::vector<std::string> get_loggers();
 void add_sink(const std::string & name, std::shared_ptr<spdlog::sinks::sink> sink);
 std::map<std::string, std::shared_ptr<spdlog::sinks::sink>> initial_sinks() {
  return sinks_;
 }
 
 static const char* appender_prefix;
 static const char* logger_prefix;
private:
 std::map<std::string, std::shared_ptr<spdlog::sinks::sink>> sinks_;
};

class LoggerConfiguration {
 public:
  static std::shared_ptr<LoggerConfiguration> getConfiguration() {
   return configuration_;
  }
  static void initialize(const std::shared_ptr<LoggerProperties> & logger_properties) {
   configuration_ = std::shared_ptr<LoggerConfiguration>(new LoggerConfiguration(logger_properties));
  }
  std::shared_ptr<spdlog::logger> get_logger(const std::string & name);
 private:
  LoggerConfiguration(const std::shared_ptr<LoggerProperties>  & logger_properties);
  static std::shared_ptr<LoggerConfiguration> configuration_;
  std::shared_ptr<LoggerNamespace> root_namespace;
};

template<typename T>
class LoggerFactory {
 public:
  static Logger& getLogger() {
   static LoggerImpl logger;
   return logger;
  }
 private:
  class LoggerImpl : public Logger {
    public:
     LoggerImpl():Logger(LoggerConfiguration::getConfiguration()->get_logger(core::getClassName<T>())) {
   }
  };
};

} /* namespace logging */
} /* namespace core */
} /* namespace minifi */
} /* namespace nifi */
} /* namespace apache */
} /* namespace org */

#endif
