/**
 * @file Logger.h
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
#ifndef __LOGGER_H__
#define __LOGGER_H__

#include <memory>

#include "spdlog/spdlog.h"
#include "spdlog/fmt/bundled/ostream.h"

namespace org {
namespace apache {
namespace nifi {
namespace minifi {
namespace core {
namespace logging {

class Logger {
 public:
  /**
   * @brief Log error message
   * @param format format string ('man printf' for syntax)
   * @warning does not check @p log or @p format for null. Caller must ensure parameters and format string lengths match
   */
  template<typename ... Args>
  void log_error(const char * const format, const Args& ... args) {
   delegate_->log(spdlog::level::err, format, args...);
  }
  
  /**
   * @brief Log warn message
   * @param format format string ('man printf' for syntax)
   * @warning does not check @p log or @p format for null. Caller must ensure parameters and format string lengths match
   */
  template<typename ... Args>
  void log_warn(const char * const format, const Args& ... args) {
   delegate_->log(spdlog::level::warn, format, args...);
  }
  
  /**
   * @brief Log info message
   * @param format format string ('man printf' for syntax)
   * @warning does not check @p log or @p format for null. Caller must ensure parameters and format string lengths match
   */
  template<typename ... Args>
  void log_info(const char * const format, const Args& ... args) {
   delegate_->log(spdlog::level::info, format, args...);
  }
  
  /**
   * @brief Log debug message
   * @param format format string ('man printf' for syntax)
   * @warning does not check @p log or @p format for null. Caller must ensure parameters and format string lengths match
   */
  template<typename ... Args>
  void log_debug(const char * const format, const Args& ... args) {
   delegate_->log(spdlog::level::debug, format, args...);
  }
  
  /**
   * @brief Log trace message
   * @param format format string ('man printf' for syntax)
   * @warning does not check @p log or @p format for null. Caller must ensure parameters and format string lengths match
   */
  template<typename ... Args>
  void log_trace(const char * const format, const Args& ... args) {
   delegate_->log(spdlog::level::trace, format, args...);
  }
  
 protected:
  Logger(std::shared_ptr<spdlog::logger> delegate) {
   delegate_ = delegate;
  }
 private:
  Logger(Logger const&);
  Logger& operator=(Logger const&);
  
  std::shared_ptr<spdlog::logger> delegate_;
};

} /* namespace logging */
} /* namespace core */
} /* namespace minifi */
} /* namespace nifi */
} /* namespace apache */
} /* namespace org */

#endif
