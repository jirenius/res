package test

import (
	"encoding/json"
	"testing"

	"github.com/jirenius/go-res"
)

// Test auth response with result
func TestAuth(t *testing.T) {
	result := `{"foo":"bar","zoo":42}`

	runTest(t, func(s *Session) {
		s.Handle("model", res.Auth("method", func(r res.AuthRequest) {
			r.OK(json.RawMessage(result))
		}))
	}, func(s *Session) {
		inb := s.Request("auth.test.model.method", nil)
		s.GetMsg(t).Equals(t, inb, json.RawMessage(`{"result":`+result+`}`))
	})
}

// Test auth response with nil result
func TestAuthWithNil(t *testing.T) {
	runTest(t, func(s *Session) {
		s.Handle("model", res.Auth("method", func(r res.AuthRequest) {
			r.OK(nil)
		}))
	}, func(s *Session) {
		inb := s.Request("auth.test.model.method", nil)
		s.GetMsg(t).Equals(t, inb, json.RawMessage(`{"result":null}`))
	})
}

// Test calling NotFound on a auth request results in system.notFound
func TestAuthNotFound(t *testing.T) {
	runTest(t, func(s *Session) {
		s.Handle("model", res.Auth("method", func(r res.AuthRequest) {
			r.NotFound()
		}))
	}, func(s *Session) {
		inb := s.Request("auth.test.model.method", nil)
		s.GetMsg(t).
			AssertSubject(t, inb).
			AssertError(t, res.ErrNotFound)
	})
}

// Test calling MethodNotFound on a auth request results in system.methodNotFound
func TestAuthMethodNotFound(t *testing.T) {
	runTest(t, func(s *Session) {
		s.Handle("model", res.Auth("method", func(r res.AuthRequest) {
			r.MethodNotFound()
		}))
	}, func(s *Session) {
		inb := s.Request("auth.test.model.method", nil)
		s.GetMsg(t).
			AssertSubject(t, inb).
			AssertError(t, res.ErrMethodNotFound)
	})
}

// Test calling InvalidParams with no message on a auth request results in system.invalidParams
func TestAuthDefaultInvalidParams(t *testing.T) {
	runTest(t, func(s *Session) {
		s.Handle("model", res.Auth("method", func(r res.AuthRequest) {
			r.InvalidParams("")
		}))
	}, func(s *Session) {
		inb := s.Request("auth.test.model.method", nil)
		s.GetMsg(t).
			AssertSubject(t, inb).
			AssertError(t, res.ErrInvalidParams)
	})
}

// Test calling InvalidParams on a auth request results in system.invalidParams
func TestAuthInvalidParams(t *testing.T) {
	runTest(t, func(s *Session) {
		s.Handle("model", res.Auth("method", func(r res.AuthRequest) {
			r.InvalidParams("foo")
		}))
	}, func(s *Session) {
		inb := s.Request("auth.test.model.method", nil)
		s.GetMsg(t).
			AssertSubject(t, inb).
			AssertError(t, &res.Error{
				Code:    res.CodeInvalidParams,
				Message: "foo",
			})
	})
}

// Test calling Error on a auth request results in given error
func TestAuthError(t *testing.T) {
	runTest(t, func(s *Session) {
		s.Handle("model", res.Auth("method", func(r res.AuthRequest) {
			r.Error(res.ErrDisposing)
		}))
	}, func(s *Session) {
		inb := s.Request("auth.test.model.method", nil)
		s.GetMsg(t).
			AssertSubject(t, inb).
			AssertError(t, res.ErrDisposing)
	})
}

// Test calling RawParams on a auth request with parameters
func TestAuthRawParams(t *testing.T) {
	params := json.RawMessage(`{"foo":"bar","baz":42}`)

	runTest(t, func(s *Session) {
		s.Handle("model", res.Auth("method", func(r res.AuthRequest) {
			AssertEqual(t, "RawParams", r.RawParams(), params)
			r.NotFound()
		}))
	}, func(s *Session) {
		req := newDefaultRequest()
		req.Params = params
		inb := s.Request("auth.test.model.method", req)
		s.GetMsg(t).
			AssertSubject(t, inb).
			AssertError(t, res.ErrNotFound)
	})
}

// Test calling RawParams on a auth request with no parameters
func TestAuthRawParamsWithNilParams(t *testing.T) {
	runTest(t, func(s *Session) {
		s.Handle("model", res.Auth("method", func(r res.AuthRequest) {
			AssertEqual(t, "RawParams", r.RawParams(), nil)
			r.NotFound()
		}))
	}, func(s *Session) {
		req := newDefaultRequest()
		inb := s.Request("auth.test.model.method", req)
		s.GetMsg(t).
			AssertSubject(t, inb).
			AssertError(t, res.ErrNotFound)
	})
}

// Test calling RawToken on a auth request with token
func TestAuthRawToken(t *testing.T) {
	token := json.RawMessage(`{"user":"foo","id":42}`)

	runTest(t, func(s *Session) {
		s.Handle("model", res.Auth("method", func(r res.AuthRequest) {
			AssertEqual(t, "RawToken", r.RawToken(), token)
			r.NotFound()
		}))
	}, func(s *Session) {
		req := newDefaultRequest()
		req.Token = token
		inb := s.Request("auth.test.model.method", req)
		s.GetMsg(t).
			AssertSubject(t, inb).
			AssertError(t, res.ErrNotFound)
	})
}

// Test calling RawToken on a auth request with no token
func TestAuthRawTokenWithNoToken(t *testing.T) {
	runTest(t, func(s *Session) {
		s.Handle("model", res.Auth("method", func(r res.AuthRequest) {
			AssertEqual(t, "RawToken", r.RawToken(), nil)
			r.NotFound()
		}))
	}, func(s *Session) {
		req := newDefaultRequest()
		inb := s.Request("auth.test.model.method", req)
		s.GetMsg(t).
			AssertSubject(t, inb).
			AssertError(t, res.ErrNotFound)
	})
}

// Test calling ParseParams on a auth request with parameters
func TestAuthParseParams(t *testing.T) {
	params := json.RawMessage(`{"foo":"bar","baz":42}`)
	var p struct {
		Foo string `json:"foo"`
		Baz int    `json:"baz"`
	}

	runTest(t, func(s *Session) {
		s.Handle("model", res.Auth("method", func(r res.AuthRequest) {
			r.ParseParams(&p)
			AssertEqual(t, "p.Foo", p.Foo, "bar")
			AssertEqual(t, "p.Baz", p.Baz, 42)
			r.NotFound()
		}))
	}, func(s *Session) {
		req := newDefaultRequest()
		req.Params = params
		inb := s.Request("auth.test.model.method", req)
		s.GetMsg(t).
			AssertSubject(t, inb).
			AssertError(t, res.ErrNotFound)
	})
}

// Test calling ParseParams on a auth request with no parameters
func TestAuthParseParamsWithNilParams(t *testing.T) {
	var p struct {
		Foo string `json:"foo"`
		Baz int    `json:"baz"`
	}

	runTest(t, func(s *Session) {
		s.Handle("model", res.Auth("method", func(r res.AuthRequest) {
			r.ParseParams(&p)
			AssertEqual(t, "p.Foo", p.Foo, "")
			AssertEqual(t, "p.Baz", p.Baz, 0)
			r.NotFound()
		}))
	}, func(s *Session) {
		req := newDefaultRequest()
		inb := s.Request("auth.test.model.method", req)
		s.GetMsg(t).
			AssertSubject(t, inb).
			AssertError(t, res.ErrNotFound)
	})
}

// Test calling ParseToken on a auth request with token
func TestAuthParseToken(t *testing.T) {
	token := json.RawMessage(`{"user":"foo","id":42}`)
	var o struct {
		User string `json:"user"`
		ID   int    `json:"id"`
	}

	runTest(t, func(s *Session) {
		s.Handle("model", res.Auth("method", func(r res.AuthRequest) {
			r.ParseToken(&o)
			AssertEqual(t, "o.User", o.User, "foo")
			AssertEqual(t, "o.ID", o.ID, 42)
			r.NotFound()
		}))
	}, func(s *Session) {
		req := newDefaultRequest()
		req.Token = token
		inb := s.Request("auth.test.model.method", req)
		s.GetMsg(t).
			AssertSubject(t, inb).
			AssertError(t, res.ErrNotFound)
	})
}

// Test calling ParseToken on a auth request with no token
func TestAuthParseTokenWithNilToken(t *testing.T) {
	var o struct {
		User string `json:"user"`
		ID   int    `json:"id"`
	}

	runTest(t, func(s *Session) {
		s.Handle("model", res.Auth("method", func(r res.AuthRequest) {
			r.ParseToken(&o)
			AssertEqual(t, "o.User", o.User, "")
			AssertEqual(t, "o.ID", o.ID, 0)
			r.NotFound()
		}))
	}, func(s *Session) {
		req := newDefaultRequest()
		inb := s.Request("auth.test.model.method", req)
		s.GetMsg(t).
			AssertSubject(t, inb).
			AssertError(t, res.ErrNotFound)
	})
}

// Test that registering auth methods with duplicate names causes panic
func TestRegisteringDuplicateAuthMethodPanics(t *testing.T) {
	runTest(t, func(s *Session) {
		defer func() {
			v := recover()
			if v == nil {
				t.Errorf(`expected test to panic, but nothing happened`)
			}
		}()
		s.Handle("model",
			res.Auth("foo", func(r res.AuthRequest) {
				r.OK(nil)
			}),
			res.Auth("bar", func(r res.AuthRequest) {
				r.OK(nil)
			}),
			res.Auth("foo", func(r res.AuthRequest) {
				r.OK(nil)
			}),
		)
	}, nil)
}
