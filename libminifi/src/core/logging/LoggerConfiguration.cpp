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
#include <string>

#include "utils/StringUtils.h"
#include "spdlog/spdlog.h"
#include "spdlog/sinks/stdout_sinks.h"
#include "spdlog/sinks/null_sink.h"

namespace org {
namespace apache {
namespace nifi {
namespace minifi {
namespace core {
namespace logging {
  
std::vector< std::string > LoggerProperties::get_appenders() {
  std::vector<std::string> appenders;
  return appenders;
}

std::vector< std::string > LoggerProperties::get_loggers() {
  std::vector<std::string> loggers;
  return loggers;
}
  
std::shared_ptr<LoggerConfiguration> LoggerConfiguration::configuration_;

LoggerConfiguration::LoggerConfiguration (LoggerProperties* logger_properties) {
  std::map<std::string, std::shared_ptr<spdlog::sinks::sink>> sink_map;
  
  for (auto const & appender_name : logger_properties->get_appenders()) {
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
  
  for (auto const & logger_key : logger_properties->get_loggers()) {
    std::string logger_def;
    if (!logger_properties->get(logger_key, logger_def)) {
      continue;
    }
    bool first = true;
    spdlog::level::level_enum level = spdlog::level::info;
    std::vector<std::shared_ptr<spdlog::sinks::sink>> sinks;
    for (auto const & segment : utils::StringUtils::split(logger_def, ",")) {
      std::string level_name = utils::StringUtils::trim(segment);
      if (first) {
        first = false;
        std::transform(level_name.begin(), level_name.end(), level_name.begin(), ::tolower);
        if ("trace" == level_name) {
          level = spdlog::level::trace;
        } else if ("debug" == level_name) {
          level = spdlog::level::debug;
        } else if ("warn" == level_name) {
          level = spdlog::level::warn;
        } else if ("critical" == level_name) {
          level = spdlog::level::critical;
        } else if ("error" == level_name) {
          level = spdlog::level::err;
        } else if ("off" == level_name) {
          level = spdlog::level::off;
        }
      } else {
        sinks.push_back(sink_map[level_name]);
      }
    }
    std::shared_ptr<LoggerNamespace> current_namespace = root_namespace;
    std::string prefix = "logger.";
    if (logger_key != "logger.root") {
      for (auto const & name : utils::StringUtils::split(logger_key.substr(prefix.length(), logger_key.length() - prefix.length()), "::")) {
        auto child_pair = current_namespace->children.find(name);
        std::shared_ptr<LoggerNamespace> child;
        if (child_pair == current_namespace->children.end()) {
          child = std::make_shared<LoggerNamespace>();
          current_namespace->children[name] = child;
        } else {
          child = child_pair->second;
        }
        current_namespace = child;
      }
    }
    current_namespace->level = level;
    current_namespace->sinks = sinks;
  }
//   std::vector<std::string> logger_properties->get_loggers();
};

std::shared_ptr<spdlog::logger> LoggerConfiguration::get_logger (const std::string& name) {
  std::shared_ptr<spdlog::logger> logger = spdlog::get(name);
  if (logger) {
    return logger;
  }
  std::shared_ptr<LoggerNamespace> current_namespace = root_namespace;
  for (auto const & name_segment : utils::StringUtils::split(name, "::")) {
    auto child_pair = current_namespace->children.find(name_segment);
    if (child_pair == current_namespace->children.end()) {
      break;
    }
    current_namespace = child_pair->second;
  }
  logger = std::make_shared<spdlog::logger>("mylogger", begin(current_namespace->sinks), end(current_namespace->sinks));
  logger->set_level(current_namespace->level);
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
