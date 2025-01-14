import 'package:flutter/material.dart';
import 'package:auto_size_text/auto_size_text.dart';
import 'dart:convert';
import "package:http/src/response.dart";
import "package:namer_app/apiCall/apiData.dart";
import 'dart:developer' as developer;

import "../custom_widget.dart";
import "sign_up.dart";
import "../explore/explorePage.dart";
import "../globalData.dart";
import "loadData/loadServicesData.dart";
import '../apiCall/apiRequest.dart';
import "loginNavigator.dart";
import "../auth/webviewAuth/connectServiceAuth.dart";

class LoginPage extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      theme: new ThemeData(
        scaffoldBackgroundColor: const Color.fromARGB(255, 255, 255, 255),
        dialogBackgroundColor: const Color.fromARGB(255, 255, 255, 255)
      ),
      home: APICall(),
    );
  }
}

class APICall extends StatefulWidget {
  @override
  _APICallState createState() => _APICallState();
}

class _APICallState extends State<APICall> {
  final TextEditingController emailController = new TextEditingController(text: "");
  final TextEditingController passwordController = new TextEditingController(text: "");
  String result = '';
  String apiRoute = 'login';

  @override
  void dispose() {
    emailController.dispose();
    passwordController.dispose();
    super.dispose();
  }

  @override
  void initState() {
    super.initState();
    _checkJWToken();
  }

  void _checkJWToken() {
    if (Globaldata.JWToken.isNotEmpty) {
      WidgetsBinding.instance.addPostFrameCallback((_) {
        Navigator.pushReplacement(
          context,
          MaterialPageRoute(builder: (context) => ExplorePage(),),
        );
      });
    }
  }

  Future<void> saveJWToken(Response response) async {
    String? cookiesBrute = response.headers["set-cookie"];
    String JWToken = "";
    if (cookiesBrute != null) {
      JWToken = cookiesBrute.split(';').first;
      setState(() {
        result = "";
        developer.log("JWToken : {${JWToken}}");
        Globaldata.JWToken = JWToken;
        LoginNavigator.redirectAndHandleGoBackRedirection(LoadServicePage(), LoadServicePage());
      });
      return;
    }
    setState(() {
      result = 'Failed to get the token';
    });
  }

  Future<void> postLoginData() async {
    final response = await ApiRequest.post(apiRoute, jsonEncode(<String, dynamic>{'email': emailController.text,'password': passwordController.text,}),);
    if (response.statusCode == 504) {
      setState(() {
        result = response.body;
      });
      return;
    }
    final responseData = jsonDecode(response.body);
    if (response.statusCode != 200) {
      developer.log(
          "postLoginData FAIL ERROR : \nError Message : $responseData\nError Code : ${response.statusCode}");
      setState(() {
        result = '${responseData['error']}';
      });
      return;
    }
    saveJWToken(response);
  }

  Future<void> changeIpDialog() async {
    final TextEditingController ipController = new TextEditingController(text: ApiData.apiUrl);
    return showDialog<void>(
      context: Globaldata.myContext,
      barrierDismissible: false,
      builder: (BuildContext context) {
        return AlertDialog(
          content: SingleChildScrollView(
            child: ListBody(
              children: <Widget>[
                Text("Select new backend ip",
                  style: TextStyle(
                    fontSize: 30,
                    fontWeight: FontWeight.w900,
                    color: Colors.black,
                  ),
                  textAlign: TextAlign.center,
                ),
                CustomClassicForm(
                  hint: 'ip : ',
                  controller: ipController,
                  maxLines: 2,
                ),
              ],
            ),
          ),
          actions: <Widget>[
            CustomButton(
              backgroundColor: const Color.fromARGB(255, 34, 34, 34),
              text: Text("Save",style: TextStyle(fontSize: 20, fontWeight: FontWeight.bold, color: Colors.white),),
              onPressed: () {
                setState(() {
                  ApiData.apiUrl = ipController.text;
                });
                Navigator.of(context).pop();
              },
              size: const Size(double.infinity, 55),
              padding: const EdgeInsets.symmetric(horizontal: 15.0),
            ),
            Center(
              child: CustomTextButton(
                text: Text("Cancel",style: TextStyle(fontSize: 20, fontWeight: FontWeight.bold, color: Colors.black),),
                onPressed: () {Navigator.of(context).pop();},
              ),
            ),
          ],
        );
      },
    );
  }

  @override
  Widget build(BuildContext context) {
    Globaldata.myContext = context;
    return Scaffold(
      body: Center(
        child: SingleChildScrollView(
          child: Column(
            children: [
              Padding(padding: EdgeInsets.only(top: 30)),
              Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  EditWheel(
                    color: Colors.black,
                    onPressed: changeIpDialog,
                  ),
                  Text("AREA",style: TextStyle(fontSize: 25, fontWeight: FontWeight.bold),),
                  Padding(padding: EdgeInsets.only(left: 48)),
                ],
              ),
              Text(
                "Log in",
                style: TextStyle(fontSize: 40, fontWeight: FontWeight.bold),
              ),
              Text(
                result,
                style: TextStyle(
                    fontSize: 20,
                    fontWeight: FontWeight.bold,
                    color: Colors.red),
              ),
              CustomClassicForm(
                hint: 'Email',
                controller: emailController,
              ),
              CustomClassicForm(
                  hint: 'Password',
                  controller: passwordController,
                  obscureText: true),
              SizedBox(height: 20),
              CustomButton(
                backgroundColor: Color.fromARGB(255, 34, 34, 34),
                text: AutoSizeText(
                  "Log in",
                  style: TextStyle(
                      fontSize: 30,
                      fontWeight: FontWeight.bold,
                      color: Colors.white),
                ),
                onPressed: postLoginData,
                size: Size(double.infinity, 70),
                padding: EdgeInsets.symmetric(horizontal: 15.0),
              ),
              Padding(padding: EdgeInsets.only(top: 15)),
              MyDivider(
                text: AutoSizeText(
                  "or",
                  style: TextStyle(fontSize: 16, color: Colors.grey),
                  maxLines: 1,
                ),
                textPadding: EdgeInsets.symmetric(horizontal: 8.0),
              ),
              Padding(padding: EdgeInsets.only(top: 5)),
              AuthButton(
                text: "Continue with Spotify",
                color: Color.fromARGB(255, 30, 215, 96),
                imageUrl: 'assets/icon/Spotify.webp',
                onPressed: () {ConnectServiceAuthWebview.connectService("login", "Spotify");},
                imageHeight: 30,
                imageWidth: 30,
              ),
              AuthButton(
                text: "Continue with Discord",
                color: Color(0xff7289da),
                imageUrl: 'assets/icon/Discord.webp',
                onPressed: () {ConnectServiceAuthWebview.connectService("login", "Discord");},
                imageHeight: 30,
                imageWidth: 30,
              ),
              AuthButton(
                text: "Continue with Github",
                color: Color(0xff24292e),
                imageUrl: 'assets/icon/Github.webp',
                onPressed: () {ConnectServiceAuthWebview.connectService("login", "Github");},
                imageHeight: 30,
                imageWidth: 30,
              ),
              AuthButton(
                text: "Continue with Gitlab",
                color: Color(0xfffc6d27),
                imageUrl: 'assets/icon/Gitlab.webp',
                onPressed: () {ConnectServiceAuthWebview.connectService("login", "Gitlab");},
                imageHeight: 30,
                imageWidth: 30,
              ),
              Padding(padding: EdgeInsets.only(top: 30)),
              CustomTextButton(
                text: Text(
                  "New to AREA? Sign up here.",
                  style: TextStyle(
                    fontSize: 15,
                    fontWeight: FontWeight.bold,
                    color: Colors.black,
                  ),
                ),
                onPressed: () {
                  Navigator.push(
                    Globaldata.myContext,
                    MaterialPageRoute(builder: (context) => SignUp()),
                  );
                },
              ),
            ],
          ),
        ),
      ),
    );
  }
}

class AuthButton extends StatelessWidget {
  const AuthButton({
    super.key,
    required this.text,
    this.color = Colors.black,
    this.imageUrl = "assets/icon/NoImage.webp",
    this.imageWidth = 24,
    this.imageHeight = 24,
    required this.onPressed,
  });

  final String text;
  final Color color;
  final String imageUrl;
  final double imageWidth;
  final double imageHeight;
  final void Function() onPressed;

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: EdgeInsets.only(top: 10),
      child: CustomButton(
        backgroundColor: color,
        text: AutoSizeText(
          text,
          style: TextStyle(
              fontSize: 25,
              fontWeight: FontWeight.bold,
              color: Colors.white),
          maxLines: 1,
        ),
        onPressed: onPressed,
        size: Size(double.infinity, 60),
        padding: EdgeInsets.symmetric(horizontal: 15.0),
        textPaddingLeft: 30,
        image: Image.asset(imageUrl,
            width: imageWidth, height: imageHeight),
      ),
    );
  }
}
