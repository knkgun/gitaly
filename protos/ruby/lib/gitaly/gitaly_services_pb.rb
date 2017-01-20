# Generated by the protocol buffer compiler.  DO NOT EDIT!
# Source: gitaly.proto for package 'gitaly'

require 'grpc'
require 'gitaly_pb'

module Gitaly
  module SmartHttp
    # The Git 'smart HTTP' protocol
    class Service

      include GRPC::GenericService

      self.marshal_class_method = :encode
      self.unmarshal_class_method = :decode
      self.service_name = 'gitaly.SmartHttp'

      # The response body for GET /info/refs?service=git-upload-pack
      rpc :InfoRefsUploadPack, InfoRefsRequest, stream(InfoRefsResponse)
      # The response body for GET /info/refs?service=git-receive-pack
      rpc :InfoRefsReceivePack, InfoRefsRequest, stream(InfoRefsResponse)
    end

    Stub = Service.rpc_stub_class
  end
end
