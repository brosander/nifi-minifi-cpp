/**
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
#ifndef LIBMINIFI_INCLUDE_UTILS_ID_H_
#define LIBMINIFI_INCLUDE_UTILS_ID_H_

#include <atomic>
#include <memory>
#include <string>
#include <uuid/uuid.h>

#include "core/logging/Logger.h"
#include "core/Property.h"

namespace org {
namespace apache {
namespace nifi {
namespace minifi {
namespace utils {

class IdGenerator {
 public:
  IdGenerator();
  void generate(uuid_t output);
  void initialize(const std::shared_ptr<Properties> & properties);
  
  static std::shared_ptr<IdGenerator> getIdGenerator() {
    static std::shared_ptr<IdGenerator> generator = std::make_shared<IdGenerator>();
    return generator;
  }
 protected:
  uint64_t getDeviceSegmentFromString(const std::string & str, int numBits);
  uint64_t getRandomDeviceSegment(int numBits);
 private:
  std::shared_ptr<logging::Logger> logger_;
  bool deterministic_;
  unsigned char deterministic_prefix_[8];
  std::atomic<uint64_t> incrementor_;
  static uint64_t char_mask;
};

class NonRepeatingStringGenerator {
 public:
  std::string generate() {
    return prefix_ + std::to_string(incrementor_++);
  }
 private:
  std::atomic<uint64_t> incrementor_;
  static std::string prefix_;
};

} /* namespace utils */
} /* namespace minifi */
} /* namespace nifi */
} /* namespace apache */
} /* namespace org */

#endif /* LIBMINIFI_INCLUDE_UTILS_ID_H_ */
