/**
 * @file Logger.cpp
 * Logger class implementation
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
#include "core/logging/LoggerConfiguration.h"
#include "utils/StringUtils.h"

#include <algorithm>
#include <vector>
#include <queue>
#include <memory>
#include <map>

#include "spdlog/spdlog.h"
#include "spdlog/sinks/stdout_sinks.h"
#include "spdlog/sinks/null_sink.h"

namespace org {
namespace apache {
namespace nifi {
namespace minifi {
namespace core {
namespace logging {

LoggerConfiguration::LoggerConfiguration (LoggerProperties* logger_properties) {
  std::vector<std::string> appender_names/* = logger_properties->get_appenders()*/;
  
  for (auto const & appender_name : appender_names) {
    std::string appender_type;
    if (!logger_properties->get(appender_name + ".type", appender_type)) {
      appender_type = "stderr";
    }
    std::transform(appender_type.begin(), appender_type.end(), appender_type.begin(), ::tolower);

    if ("nullappender" == appender_type || "null appender" == appender_type || "null" == appender_type) {
      sink_map[appender_name] = std::make_shared<spdlog::sinks::null_sink_st>();
    } else if ("rollingappender" == appender_type || "rolling appender" == appender_type || "rolling" == appender_type) {
      std::string file_name = "";
      if (!logger_properties->get(appender_name + ".file_name", file_name)) {
        file_name = "minifi-app";
      }
      std::string file_ext = "";
      if (!logger_properties->get(appender_name + ".file_extension", file_ext)) {
        file_ext = "log";
      }

      int max_files = 3;
      std::string max_files_str = "";
      if (logger_properties->get(appender_name + ".max_files", max_files_str)) {
        try {
          max_files = std::stoi(max_files_str);
        } catch (const std::invalid_argument &ia) {
        } catch (const std::out_of_range &oor) {}
      }

      int max_file_size = 5*1024*1024;
      std::string max_file_size_str = "";
      if (logger_properties->get(appender_name + ".max_file_size", max_file_size_str)) {
        try {
          max_file_size = std::stoi(max_file_size_str);
        } catch (const std::invalid_argument &ia) {
        } catch (const std::out_of_range &oor) {}
      }
      sink_map[appender_name] = std::make_shared<spdlog::sinks::rotating_file_sink_mt>(file_name, file_ext, max_file_size, max_files);
    } else if ("outputstream" == appender_name || "outputstreamappender" == appender_name || "outputstream appender" == appender_name) {
      sink_map[appender_name] = spdlog::sinks::stdout_sink_mt::instance();
    } else {
      sink_map[appender_name] = spdlog::sinks::stderr_sink_mt::instance();
    }
  }
  sink_map["root"] = spdlog::sinks::stdout_sink_mt::instance();
  
//   std::vector<std::string> logger_properties->get_loggers();
};

std::shared_ptr<spdlog::logger> LoggerConfiguration::get_logger (const std::__cxx11::string& name) {
  std::shared_ptr<spdlog::logger> logger = spdlog::get(name);
  if (logger) {
    return logger;
  }
  std::shared_ptr<spdlog::sinks::sink> sink = sink_map.find("root")->second;
  logger = std::make_shared<spdlog::logger>("mylogger", sink);
  logger->set_level(spdlog::level::info);
  try {
    spdlog::register_logger(logger);
  } catch (const spdlog::spdlog_ex &ex) {
    // Ignore as someone else beat us to registration, we should get the one they made below
  }
  return spdlog::get(name);
}


} /* namespace logging */
} /* namespace core */
} /* namespace minifi */
} /* namespace nifi */
} /* namespace apache */
} /* namespace org */
