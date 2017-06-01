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

#include <uuid/uuid.h>

#include "core/Id.h"

std::string org::apache::nifi::minifi::core::NonRepeatingStringGenerator::prefix_ =
  std::to_string(std::chrono::duration_cast<std::chrono::milliseconds>(std::chrono::system_clock::now().time_since_epoch()).count()) + "-";

org::apache::nifi::minifi::core::Id::Id() : initialized_(false) {}

org::apache::nifi::minifi::core::Id::Id(const std::string& uuid_str) : initialized_(true) {
    uuid_parse(uuid_str.c_str(), uuid_);
    char uuidStr[37];
    uuid_unparse_lower(uuid_, uuidStr);
    uuidStr_ = uuidStr;
}

org::apache::nifi::minifi::core::Id::Id(uuid_t uuid) {
  if (uuid) {
    uuid_copy(uuid_, uuid);
    char uuidStr[37];
    uuid_unparse_lower(uuid_, uuidStr);
    uuidStr_ = uuidStr;
    initialized_ = true;
  } else {
    initialized_ = false;
  }
}

void org::apache::nifi::minifi::core::Id::initializeIfNecessary() {
  std::lock_guard<std::mutex> lock(mutex_);
  if (!initialized_) {
    uuid_generate(uuid_);
    char uuidStr[37];
    uuid_unparse_lower(uuid_, uuidStr);
    uuidStr_ = uuidStr;
    initialized_ = true;
  }
}

// Get UUID
bool org::apache::nifi::minifi::core::Id::getUUID(uuid_t uuid) {
  initializeIfNecessary();
  if (uuid) {
    uuid_copy(uuid, uuid_);
    return true;
  } else {
    return false;
  }
}

const std::string& org::apache::nifi::minifi::core::Id::getUUIDStr() {
  initializeIfNecessary();
  return uuidStr_;
}

// Get UUID
unsigned const char *org::apache::nifi::minifi::core::Id::getUUID() {
  initializeIfNecessary();
  return uuid_;
}