# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: teststream.proto

require 'google/protobuf'

require 'lint_pb'
require 'shared_pb'
require 'google/protobuf/empty_pb'
Google::Protobuf::DescriptorPool.generated_pool.build do
  add_file("teststream.proto", :syntax => :proto3) do
    add_message "gitaly.TestStreamRequest" do
      optional :repository, :message, 1, "gitaly.Repository"
      optional :size, :int64, 2
    end
  end
end

module Gitaly
  TestStreamRequest = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("gitaly.TestStreamRequest").msgclass
end
