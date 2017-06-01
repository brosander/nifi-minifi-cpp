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
#ifndef LIBMINIFI_INCLUDE_CORE_ID_H_
#define LIBMINIFI_INCLUDE_CORE_ID_H_

#include <atomic>
#include <mutex>
#include <string>
#include <uuid/uuid.h>

namespace org {
namespace apache {
namespace nifi {
namespace minifi {
namespace core {

class NonRepeatingStringGenerator {
 public:
  std::string generate() {
    return prefix_ + std::to_string(incrementor_++);
  }
 private:
  std::atomic<uint64_t> incrementor_;
  static std::string prefix_;
};

class Id {
 public:
  Id();
  Id(const std::string &uuid_str);
  Id(uuid_t uuid);

  /**
   * Set UUID in this instance
   * @param uuid uuid to apply to the internal representation.
   */
  void setUUID(uuid_t uuid);

  /**
   * Set UUID in this instance
   * @param uuid uuid to apply to the internal representation.
   */
  void setUUID(const std::string uuid_str);

  /**
   * Returns the UUID through the provided object.
   * @param uuid uuid struct to which we will copy the memory
   * @return success of request
   */
  bool getUUID(uuid_t uuid);

  unsigned const char *getUUID();
  /**
   * Return the UUID string
   * @param constant reference to the UUID str
   */
  const std::string & getUUIDStr();
 private:
  bool initialized_;
  // A global unique identifier
  uuid_t uuid_;
  // UUID string
  std::string uuidStr_;

  std::mutex mutex_;

  inline void initializeIfNecessary();
};

} /* namespace core */
} /* namespace minifi */
} /* namespace nifi */
} /* namespace apache */
} /* namespace org */
#endif /* LIBMINIFI_INCLUDE_CORE_ID_H_ */
